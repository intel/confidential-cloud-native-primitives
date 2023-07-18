/* SPDX-License-Identifier: Apache-2.0 */

use clap::Parser;
use quote_server::get_quote_server::{GetQuote, GetQuoteServer};
use quote_server::{GetQuoteRequest, GetQuoteResponse};
use tokio::net::UnixListener;
use tokio_stream::wrappers::UnixListenerStream;
use tonic::{transport::Server, Request, Response, Status};

pub mod tee;
use tee::*;

pub mod quote_server {
    tonic::include_proto!("quoteserver");

    pub(crate) const FILE_DESCRIPTOR_SET: &[u8] =
        tonic::include_file_descriptor_set!("quote_server_descriptor");
}

pub struct CCNPGetQuote {
    local_tee: tee::TeeType,
}

impl CCNPGetQuote {
    fn new(_local_tee: TeeType) -> Self {
        CCNPGetQuote {
            local_tee: _local_tee,
        }
    }
}

#[tonic::async_trait]
impl GetQuote for CCNPGetQuote {
    async fn get_quote(
        &self,
        request: Request<GetQuoteRequest>,
    ) -> Result<Response<GetQuoteResponse>, Status> {
        let msg;
        let req = request.into_inner();

        println!(
            "Got a request with: user_data = {:?}, nonce = {:?}",
            req.user_data, req.nonce
        );
        let result = get_quote(self.local_tee.clone(), req.user_data, req.nonce);
        match result {
            Ok(q) => {
                msg = Response::new(quote_server::GetQuoteResponse {
                    quote: q,
                    quote_type: format!("{:?}", self.local_tee).to_string(),
                })
            }
            Err(e) => return Err(Status::internal(e.to_string())),
        }
        Ok(msg)
    }
}

#[derive(Parser)]
struct Cli {
    port: String,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let path = "/run/ccnp/uds/quote-server.sock";
    let _ = std::fs::remove_file(path);
    let uds = UnixListener::bind(path)?;
    let uds_stream = UnixListenerStream::new(uds);

    let getquote = CCNPGetQuote::new({
        match tee::get_tee_type() {
            tee::TeeType::PLAIN => panic!("Not found any TEE device!"),
            t => t,
        }
    });

    let (mut health_reporter, health_service) = tonic_health::server::health_reporter();
    health_reporter
        .set_serving::<GetQuoteServer<CCNPGetQuote>>()
        .await;

    let reflection_service = tonic_reflection::server::Builder::configure()
        .register_encoded_file_descriptor_set(quote_server::FILE_DESCRIPTOR_SET)
        .build()
        .unwrap();

    println!(
        "Starting quote server in {} enviroment...",
        format!("{:?}", tee::get_tee_type()).to_string()
    );

    Server::builder()
        .add_service(reflection_service)
        .add_service(health_service)
        .add_service(GetQuoteServer::new(getquote))
        .serve_with_incoming(uds_stream)
        .await?;
    Ok(())
}
