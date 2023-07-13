/* SPDX-License-Identifier: Apache-2.0 */

use anyhow::*;
use sha2::{Digest, Sha512};
use std::path::Path;
use std::result::Result::Ok;
use tdx_attest_rs;

#[derive(Debug, Clone)]
pub enum TeeType {
    TDX,
    SEV,
    TPM,
    PLAIN,
}

pub fn get_tee_type() -> TeeType {
    if Path::new("/dev/tpm0").exists() {
        TeeType::TPM
    } else if Path::new("/dev/tdx-guest").exists()
        || Path::new("/dev/tdx-attest").exists()
        || Path::new("/dev/tdx_guest").exists()
    {
        if Path::new("/dev/tdx-attest").exists() {
            panic!("get_tdx_quote: Deprecated device node /dev/tdx-attest, please upgrade to use /dev/tdx-guest or /dev/tdx_guest");
        }
        TeeType::TDX
    } else if Path::new("/dev/sev-guest").exists() || Path::new("/dev/sev").exists() {
        TeeType::SEV
    } else {
        TeeType::PLAIN
    }
}

fn generate_tdx_report_data(
    report_data: Option<String>,
    nonce: String,
) -> Result<tdx_attest_rs::tdx_report_data_t, anyhow::Error> {
    let hash = Sha512::new().chain_update(nonce.into_bytes());
    let _ret = match report_data {
        Some(_encoded_report_data) => {
            if _encoded_report_data.is_empty() {
                hash.clone()
            } else {
                let decoded_report_data = match base64::decode(_encoded_report_data) {
                    Ok(v) => v,
                    Err(e) => return Err(anyhow!("user data is not base64 encoded: {:?}", e)),
                };
                hash.clone().chain_update(decoded_report_data)
            }
        }
        None => hash.clone(),
    };
    let _d: [u8; 64] = hash
        .finalize()
        .as_slice()
        .try_into()
        .expect("Wrong length of report data");
    Ok(tdx_attest_rs::tdx_report_data_t { d: _d })
}

fn get_tdx_quote(report_data: Option<String>, nonce: String) -> Result<String> {
    if nonce.is_empty() {
        return Err(anyhow!("empty nonce!"));
    }

    let tdx_report_data = match generate_tdx_report_data(report_data, nonce) {
        Ok(v) => v,
        Err(e) => {
            return Err(anyhow!("get_tdx_quote: {:?}", e));
        }
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

pub fn get_quote(local_tee: TeeType, user_data: String, nonce: String) -> Result<String> {
    match local_tee {
        TeeType::TDX => get_tdx_quote(Some(user_data), nonce),
        TeeType::TPM => get_tpm_quote(),
        TeeType::SEV => get_sev_quote(),
        _ => Err(anyhow!("Unexpected case!")),
    }
}

#[test]
fn tdx_get_quote_no_nonce() {
    let result = get_tdx_quote(Some("YWJjZGVmZw==".to_string()), "".to_string());
    assert!(result.is_err());
}

#[test]
fn tdx_get_quote_report_data_size_8() {
    // "YWJjZGVmZw==" is base64 of "abcdefg", 8 bytes
    let result = get_tdx_quote(
        Some("YWJjZGVmZw==".to_string()),
        "IXUKoBO1XEFBPwopN4sY".to_string(),
    );
    assert!(result.is_ok());
}

#[test]
fn tdx_get_quote_report_data_size_0() {
    //allow does not specify report data
    let result = get_tdx_quote(Some("".to_string()), "IXUKoBO1XEFBPwopN4sY".to_string());
    assert!(result.is_ok());
}

#[test]
fn tdx_get_quote_report_data_size_48() {
    // this one should be standard 48 bytes base64 encoded report data
    // "MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4" is base64 of "123456781234567812345678123456781234567812345678", 48 bytes
    let result = get_tdx_quote(
        Some("MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4".to_string()),
        "IXUKoBO1XEFBPwopN4sY".to_string(),
    );
    assert!(result.is_ok());
}

#[test]
fn tdx_get_quote_report_data_null() {
    // allow call get_tdx_quote with out specify report data
    let result = get_tdx_quote(None, "IXUKoBO1XEFBPwopN4sY".to_string());
    assert!(result.is_ok());
}

#[test]
fn tdx_get_quote_report_data_not_base64_encoded() {
    //coming in report data should always be base64 encoded
    let result = get_tdx_quote(
        Some("XD^%*!x".to_string()),
        "IXUKoBO1XEFBPwopN4sY".to_string(),
    );
    assert!(result.is_err());
}

#[test]
fn tdx_get_quote_long_tdx_report_data() {
    let result = get_tdx_quote(
        Some(
            "MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2Nzgx\
            MjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEy\
            MzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIz\
            NDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0\
            NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1\
            Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2\
            NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4Cg=="
                .to_string(),
        ),
        "IXUKoBO1XEFBPwopN4sY".to_string(),
    );
    assert!(result.is_ok());
}

#[test]
fn tdx_get_quote_long_nonce() {
    let result = get_tdx_quote(
        Some("MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4".to_string()),
        "MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2Nzgx\
        MjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEy\
        MzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIz\
        NDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0\
        NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1\
        Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2\
        NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4Cg=="
            .to_string(),
    );
    assert!(result.is_ok());
}

#[test]
fn get_quote_wrong_tee_type() {
    //does not allow tee type beyond TDX/SEV/TPM
    let result = get_quote(
        TeeType::PLAIN,
        "".to_string(),
        "IXUKoBO1XEFBPwopN4sY".to_string(),
    );
    assert!(result.is_err());
}
