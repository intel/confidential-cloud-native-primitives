/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
*/

use anyhow::*;
use clap::Parser;
use core::result::Result::Ok;
use hyper::service::{make_service_fn, service_fn};
use hyper::{Request as HyperRequest, Response as HyperResponse, Body, Server as HyperServer};
use std::net::SocketAddr;

pub mod kube;
pub mod tee;
use tee::*;

// A http server for provide the current pod quote data
#[derive(Copy, Clone)]
pub struct PerPodQuoteServer {
    sock_address: SocketAddr,
    local_tee: tee::TeeType,
}

impl PerPodQuoteServer {
    pub fn new(sock_address: SocketAddr, local_tee: tee::TeeType) -> Self {
        PerPodQuoteServer {
            sock_address,
            local_tee,
        }
    }

    pub async fn start(&self) -> Result<(), hyper::Error> {
        let local_tee = self.local_tee;
        let make_svc = make_service_fn(|_conn| {
            let service = service_fn(move |req| {
                // Route request to the appropriate handler
                Self::handle_request(local_tee, req)
            });
            async move { Ok::<_, hyper::Error>(service) }
        });
        let http_server = HyperServer::bind(&self.sock_address).serve(make_svc);
        println!(
            "The Pod Quote HTTP server is listening on: {:?}",
            self.sock_address
        );
        http_server.await
    }

    // generate current pod quote based on its all containers' imageIDs
    async fn get_current_pod_quote(local_tee: tee::TeeType) -> Result<String, anyhow::Error> {
        // Handle the "/quote" route
        // Create an instance of your custom kube client
        let pod_data = kube::get_cur_pod_images_info();
        match pod_data.await {
            Ok(report_data) => {
                let report_data_clone = report_data.clone();
                let hash_report_data = kube::sha256_hash(&report_data_clone);
                let quote_data = get_quote(
                    local_tee,
                    hash_report_data.clone(),
                    hash_report_data.clone(),
                )
                .unwrap();
                Ok(quote_data)
            }
            Err(error) => Err(anyhow!(
                "There was a problem when get current pod images information: {:?}",
                error
            )),
        }
    }

    async fn handle_request(
        local_tee: tee::TeeType,
        req: HyperRequest<Body>
    ) -> Result<HyperResponse<Body>, hyper::Error> {
        match req.uri().path() {
            "/quote" => {
                match Self::get_current_pod_quote(local_tee).await {
                    Ok(quote_data) => {
                        println!("File content: {}", quote_data);
                        // generate the response from quote file
                        let response = HyperResponse::new(Body::from(quote_data));
                        Ok(response)
                    }
                    Err(err) => {
                        eprintln!("Error: {}", err);
                        let response = HyperResponse::builder()
                            .status(404)
                            .body(Body::from("Not Found Quote File"))
                            .unwrap();
                        Ok(response)
                    }
                }
            }
            _ => {
                // Handle other routes
                let response = HyperResponse::builder()
                    .status(404)
                    .body(Body::from("Not Found"))
                    .unwrap();
                Ok(response)
            }
        }
    }
}

#[derive(Parser)]
struct Cli {
    port: String,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let http_addr = SocketAddr::from(([127, 0, 0, 1], 3000));
    // Create the http server tokio task for fetching quote with current pod image IDs
    let _ = tokio::spawn(async move {
        let http_server = PerPodQuoteServer::new(http_addr, {
            match tee::get_tee_type() {
                tee::TeeType::PLAIN => panic!("Not found any TEE device!"),
                t => t,
            }
        });
        if let Err(err) = http_server.start().await {
            eprintln!("HTTP server error: {}", err);
        }
    });
    Ok(())
}
