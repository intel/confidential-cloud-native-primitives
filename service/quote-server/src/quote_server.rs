/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
*/

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
    let uds = match UnixListener::bind(path) {
        Ok(r) => r,
        Err(e) => panic!("[quote-server]: bind UDS socket error: {:?}", e),
    };
    let uds_stream = UnixListenerStream::new(uds);

    let getquote = CCNPGetQuote::new({
        match tee::get_tee_type() {
            tee::TeeType::PLAIN => panic!("[quote-server]: Not found any TEE device!"),
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

#[cfg(test)]
mod quote_server_tests {
    use super::*;
    use crate::quote_server::get_quote_client::GetQuoteClient;
    use serial_test::serial;
    use tokio::net::UnixStream;
    use tonic::transport::{Endpoint, Uri};
    use tower::service_fn;

    async fn creat_server() {
        let path = "/tmp/quote-server.sock";
        let _ = std::fs::remove_file(path);
        let uds = match UnixListener::bind(path) {
            Ok(r) => r,
            Err(e) => panic!("[quote-server]: bind UDS socket error: {:?}", e),
        };
        let uds_stream = UnixListenerStream::new(uds);

        let getquote = CCNPGetQuote::new({
            match tee::get_tee_type() {
                tee::TeeType::PLAIN => panic!("[quote-server]: Not found any TEE device!"),
                t => t,
            }
        });

        tokio::spawn(async {
            Server::builder()
                .add_service(GetQuoteServer::new(getquote))
                .serve_with_incoming(uds_stream)
                .await
                .unwrap();
        });
    }

    #[tokio::test]
    #[serial]
    //test start server and send request
    async fn request_to_server_normal() {
        creat_server().await;

        let channel = Endpoint::try_from("http://[::]:40081")
            .unwrap()
            .connect_with_connector(service_fn(|_: Uri| {
                let path = "/tmp/quote-server.sock";
                UnixStream::connect(path)
            }))
            .await
            .unwrap();

        let mut client = GetQuoteClient::new(channel);

        let request = tonic::Request::new(GetQuoteRequest {
            user_data: base64::encode("123456781234567812345678123456781234567812345678"),
            nonce: "12345678".to_string(),
        });

        let response = client.get_quote(request).await.unwrap().into_inner();
        assert_eq!(response.quote_type, "TDX");
    }

    #[tokio::test]
    #[serial]
    async fn request_to_server_empty_user_data() {
        creat_server().await;

        let channel = Endpoint::try_from("http://[::]:40081")
            .unwrap()
            .connect_with_connector(service_fn(|_: Uri| {
                let path = "/tmp/quote-server.sock";
                UnixStream::connect(path)
            }))
            .await
            .unwrap();

        let mut client = GetQuoteClient::new(channel);

        let request = tonic::Request::new(GetQuoteRequest {
            user_data: "".to_string(),
            nonce: "12345678".to_string(),
        });

        let response = client.get_quote(request).await.unwrap().into_inner();
        assert_eq!(response.quote_type, "TDX");
        assert_ne!(response.quote.len(), 0);
    }

    #[tokio::test]
    #[serial]
    async fn request_to_server_long_user_data() {
        creat_server().await;

        let channel = Endpoint::try_from("http://[::]:40081")
            .unwrap()
            .connect_with_connector(service_fn(|_: Uri| {
                let path = "/tmp/quote-server.sock";
                UnixStream::connect(path)
            }))
            .await
            .unwrap();

        let mut client = GetQuoteClient::new(channel);

        let request = tonic::Request::new(GetQuoteRequest {
            user_data: "123456781234567812345678123456781234567812345678123456781234567812345678123456781234567812345678123456781234567812345678123456781234567812345678123456781234567812345678123456781234567812345678".to_string(),
            nonce: "12345678".to_string(),
        });

        let response = client.get_quote(request).await.unwrap().into_inner();
        assert_eq!(response.quote_type, "TDX");
        assert_ne!(response.quote.len(), 0);
    }

    #[tokio::test]
    #[serial]
    async fn request_to_server_empty_nonce() {
        creat_server().await;

        let channel = Endpoint::try_from("http://[::]:40081")
            .unwrap()
            .connect_with_connector(service_fn(|_: Uri| {
                let path = "/tmp/quote-server.sock";
                UnixStream::connect(path)
            }))
            .await
            .unwrap();

        let mut client = GetQuoteClient::new(channel);

        let request = tonic::Request::new(GetQuoteRequest {
            user_data: "123456781234567812345678123456781234567812345678".to_string(),
            nonce: "".to_string(),
        });

        let response = client.get_quote(request).await.unwrap().into_inner();
        assert_eq!(response.quote_type, "TDX");
        assert_ne!(response.quote.len(), 0);
    }

    #[tokio::test]
    #[serial]
    async fn request_to_server_log_nonce() {
        creat_server().await;

        let channel = Endpoint::try_from("http://[::]:40081")
            .unwrap()
            .connect_with_connector(service_fn(|_: Uri| {
                let path = "/tmp/quote-server.sock";
                UnixStream::connect(path)
            }))
            .await
            .unwrap();

        let mut client = GetQuoteClient::new(channel);

        let request = tonic::Request::new(GetQuoteRequest {
            user_data: "123456781234567812345678123456781234567812345678".to_string(),
            nonce: "123456781234567812345678123456781234567812345678123456781234567812345678123456781234567812345678123456781234567812345678123456781234567812345678123456781234567812345678123456781234567812345678".to_string(),
        });

        let response = client.get_quote(request).await.unwrap().into_inner();
        assert_eq!(response.quote_type, "TDX");
        assert_ne!(response.quote.len(), 0);
    }

    #[tokio::test]
    #[serial]
    async fn request_to_server_verify_report_data() {
        creat_server().await;

        let channel = Endpoint::try_from("http://[::]:40081")
            .unwrap()
            .connect_with_connector(service_fn(|_: Uri| {
                let path = "/tmp/quote-server.sock";
                UnixStream::connect(path)
            }))
            .await
            .unwrap();

        let mut client = GetQuoteClient::new(channel);

        let request = tonic::Request::new(GetQuoteRequest {
            user_data: "YWJjZGVmZw==".to_string(),
            nonce: "MTIzNDU2Nzg=".to_string(),
        });

        let response = client.get_quote(request).await.unwrap().into_inner();

        let expected_report_data = [
            93, 71, 28, 83, 115, 189, 166, 130, 87, 137, 126, 119, 140, 209, 163, 215, 13, 175,
            225, 101, 64, 195, 196, 202, 15, 37, 166, 241, 141, 49, 128, 157, 164, 132, 67, 50, 9,
            32, 162, 89, 243, 191, 177, 131, 4, 159, 156, 104, 11, 193, 18, 217, 92, 215, 194, 98,
            145, 191, 211, 85, 187, 118, 39, 80,
        ];

        assert_eq!(response.quote_type, "TDX");
        let quote = base64::decode(response.quote.replace("\"", "")).unwrap();
        let mut report_data_in_quote: [u8; 64 as usize] = [0; 64 as usize];
        report_data_in_quote.copy_from_slice(&quote[568..632]);
        assert_eq!(report_data_in_quote, expected_report_data);
    }
}
