/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
*/

use std::{fs, os::unix::fs::PermissionsExt};

use anyhow::Result;
use clap::Parser;
use log::info;
use tokio::net::UnixListener;
use tokio_stream::wrappers::UnixListenerStream;
use tonic::transport::Server;
use simple_logger::SimpleLogger;

pub mod ccnp_pb {
    tonic::include_proto!("ccnp_server_pb");

    pub(crate) const FILE_DESCRIPTOR_SET: &[u8] =
        tonic::include_file_descriptor_set!("ccnp_server_descriptor");
}
use ccnp_pb::ccnp_server::CcnpServer;

mod ccnp_service;


#[derive(Parser)]
struct Cli {
    port: String,
}

const SOCK:&str = "/run/ccnp/uds/ccnp-server.sock";

fn set_sock_perm() -> Result<()> {
    let mut perms = fs::metadata(SOCK)?.permissions();
    perms.set_mode(0o666);
    fs::set_permissions(SOCK, perms)?;
    Ok(())
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    SimpleLogger::new().init()?;
    let _ = std::fs::remove_file(SOCK);
    let uds = match UnixListener::bind(SOCK) {
        Ok(r) => r,
        Err(e) => panic!("[ccnp-server]: bind UDS socket error: {:?}", e),
    };
    let uds_stream = UnixListenerStream::new(uds);
    set_sock_perm()?;

    let service = ccnp_service::Service::default();
    
    let (mut health_reporter, health_service) = tonic_health::server::health_reporter();
    health_reporter
        .set_serving::<CcnpServer<ccnp_service::Service>>()
        .await;

    let reflection_service = tonic_reflection::server::Builder::configure()
        .register_encoded_file_descriptor_set(ccnp_pb::FILE_DESCRIPTOR_SET)
        .build()
        .unwrap();

    info!("Starting ccnp server...");

    Server::builder()
        .add_service(reflection_service)
        .add_service(health_service)
        .add_service(CcnpServer::new(service))
        .serve_with_incoming(uds_stream)
        .await?;
    Ok(())
}

