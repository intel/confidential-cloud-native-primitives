/*
*
* Copyright 2023 Intel authors.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
 */

use anyhow::*;
use std::path::Path;
use tdx_attest_rs;

#[derive(Debug, Clone)]
pub enum TeeType {
    TDX,
    SEV,
    TPM,
    PLAIN,
}

pub fn get_tee_type() -> TeeType {
    if Path::new("/dev/tdx-guest").exists()
        || Path::new("/dev/tdx-attest").exists()
        || Path::new("/dev/tdx_guest").exists()
    {
        if Path::new("/dev/tdx-attest").exists() {
            panic!("get_tdx_quote: Deprecated device node /dev/tdx-attest, please upgrade to use /dev/tdx-guest or /dev/tdx_guest");
        }
        TeeType::TDX
    } else if Path::new("/dev/tpm0").exists() {
        TeeType::TPM
    } else if Path::new("/dev/sev-guest").exists() || Path::new("/dev/sev").exists() {
        TeeType::SEV
    } else {
        TeeType::PLAIN
    }
}

fn get_tdx_quote(report_data: Option<String>) -> Result<String> {
    let tdx_report_data = match report_data {
        Some(_report_data) => {
            if _report_data.is_empty() {
                tdx_attest_rs::tdx_report_data_t { d: [0u8; 64usize] }
            } else {
                let mut _tdx_report_data = base64::decode(_report_data)?;
                if _tdx_report_data.len() != 48 {
                    return Err(anyhow!(
                        "get_tdx_quote: runtime data should be SHA384 base64 String of 48 bytes"
                    ));
                }
                _tdx_report_data.extend([0; 16]);
                tdx_attest_rs::tdx_report_data_t {
                    d: _tdx_report_data.as_slice().try_into()?,
                }
            }
        }
        None => tdx_attest_rs::tdx_report_data_t { d: [0u8; 64usize] },
    };

    let quote = match tdx_attest_rs::tdx_att_get_quote(Some(&tdx_report_data), None, None, 0) {
        (tdx_attest_rs::tdx_attest_error_t::TDX_ATTEST_SUCCESS, Some(q)) => base64::encode(q),
        (error_code, _) => {
            return Err(anyhow!("get_tdx_quote: {:?}", error_code));
        }
    };

    serde_json::to_string(&quote).map_err(|e| anyhow!("get_tdx_quote: {:?}", e))
}

fn get_tpm_quote() -> Result<String> {
    Err(anyhow!("TPM to be supported!"))
}

fn get_sev_quote() -> Result<String> {
    Err(anyhow!("SEV to be supported!"))
}

pub fn get_quote(local_tee: TeeType, report_data: String) -> Result<String> {
    match local_tee {
        TeeType::TDX => get_tdx_quote(Some(report_data)),
        TeeType::TPM => get_tpm_quote(),
        TeeType::SEV => get_sev_quote(),
        _ => Err(anyhow!("Unexpected case!")),
    }
}

#[test]
fn tdx_report_data_size_8() {
    // "YWJjZGVmZw==" is base64 of "abcdefg", 8 bytes
    let result = get_tdx_quote(Some("YWJjZGVmZw==".to_string()));
    assert!(result.is_err());
}

#[test]
fn tdx_report_data_size_0() {
    //allow does not specify report data
    let result = get_tdx_quote(Some("".to_string()));
    assert!(result.is_ok());
}

#[test]
fn tdx_report_data_size_48() {
    // this one should be standard 48 bytes base64 encoded report data
    // "MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4" is base64 of "123456781234567812345678123456781234567812345678", 48 bytes
    let result = get_tdx_quote(Some(
        "MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4".to_string(),
    ));
    assert!(result.is_ok());
}

#[test]
fn tdx_report_data_null() {
    // allow call get_tdx_quote with out specify report data
    let result = get_tdx_quote(None);
    assert!(result.is_ok());
}

#[test]
fn tdx_report_data_not_base64_encoded() {
    //coming in report data should always be base64 encoded
    let result = get_tdx_quote(Some(
        "123456781234567812345678123456781234567812345678".to_string(),
    ));
    assert!(result.is_err());
}

#[test]
fn get_quote_wrong_tee_type() {
    //does not allow tee type beyond TDX/SEV/TPM
    let result = get_quote(TeeType::PLAIN, "".to_string());
    assert!(result.is_err());
}