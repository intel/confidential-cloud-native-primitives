/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
*/

extern crate crypto_hash;
extern crate kube;

use anyhow::{anyhow, Error};
use k8s_openapi::api::core::v1::Pod;
use kube::api::Api;
use kube::Client;

use std::env;

const POD_NAME: &str = "POD_NAME";
const POD_NAMESPACE: &str = "POD_NAMESPACE";
const SEPARATOR: &str = "|";

pub async fn get_cur_pod_images_info() -> Result<String, Error> {
    let mut pod_data_array: Vec<String> = Vec::new();
    let namespace = env::var(POD_NAMESPACE).unwrap_or_default();
    let pod_name = env::var(POD_NAME).unwrap_or_default();

    let client = Client::try_default().await?;
    let pods: Api<Pod> = Api::namespaced(client.clone(), &namespace);

    let pod_name_str = pod_name.clone();
    let cur_pod = pods.get(&pod_name_str).await?;
    // Access the container statuses
    if let Some(status) = cur_pod.status {
        for container_status in status.container_statuses.unwrap_or_default() {
            let image_id = container_status.image_id.clone();
            pod_data_array.push(image_id);
        }
        println!("pod quote data array:");
        // Print out the quote data of pod.
        for item in &pod_data_array {
            println!("{}", item);
        }

        // Concat all pod quote data into one String.
        let pod_image_id_data = pod_data_array.join(SEPARATOR);
        return Ok(pod_image_id_data);
    } else {
        println!("Pod {pod_name} in {namespace} not found.");
        let error_message = format!("Pod '{}' in '{}' not found.", pod_name, namespace);
        return Err(anyhow!(error_message));
    }
}

pub fn sha256_hash(input: &str) -> String {
    // Convert the input string to bytes
    let input_bytes = input.as_bytes();

    // Calculate the SHA-256 hash
    let hash = crypto_hash::hex_digest(crypto_hash::Algorithm::SHA256, input_bytes);

    hash
}
