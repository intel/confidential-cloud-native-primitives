/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
*/

use anyhow::*;
use sha2::{Digest, Sha512};
use std::path::Path;
use std::result::Result::Ok;

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
            panic!("[get_tee_type]: Deprecated device node /dev/tdx-attest, please upgrade to use /dev/tdx-guest or /dev/tdx_guest");
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
) -> Result<String, anyhow::Error> {
    let nonce_decoded = match base64::decode(nonce) {
        Ok(v) => v,
        Err(e) => {
            return Err(anyhow!(
                "[generate_tdx_report_data] nonce is not base64 encoded: {:?}",
                e
            ))
        }
    };
    let mut hasher = Sha512::new();
    hasher.update(nonce_decoded);
    let _ret = match report_data {
        Some(_encoded_report_data) => {
            if _encoded_report_data.is_empty() {
                hasher.update("")
            } else {
                let decoded_report_data = match base64::decode(_encoded_report_data) {
                    Ok(v) => v,
                    Err(e) => {
                        return Err(anyhow!(
                            "[generate_tdx_report_data] user data is not base64 encoded: {:?}",
                            e
                        ))
                    }
                };
                hasher.update(decoded_report_data)
            }
        }
        None => hasher.update(""),
    };
    let hash_array: [u8; 64] = hasher
        .finalize()
        .as_slice()
        .try_into()
        .expect("[generate_tdx_report_data] Wrong length of report data");
    Ok(base64::encode(hash_array))
}

fn get_tdx_quote(report_data: Option<String>, nonce: String) -> Result<String> {
    let tdx_report_data = match generate_tdx_report_data(report_data, nonce) {
        Ok(v) => v,
        Err(e) => {
            return Err(anyhow!("[get_tdx_quote]: {:?}", e));
        }
    };

    let quote = match tdx_attest::get_tdx_quote(tdx_report_data) {
        Err(e) => panic!("[get_tdx_quote] Fail to get TDX quote: {:?}", e),
        Ok(q) => base64::encode(q),
    };

    serde_json::to_string(&quote).map_err(|e| anyhow!("[get_tdx_quote]: {:?}", e))
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

#[cfg(test)]
mod tests {

    use super::*;

    #[test]
    //generate_tdx_report allow empty nonce
    fn generate_tdx_report_data_empty_nonce() {
        let result = generate_tdx_report_data(Some("YWJjZGVmZw==".to_string()), "".to_string());
        assert!(result.is_ok());
    }

    #[test]
    //generate_tdx_report allow optional report data
    fn tdx_get_quote_report_data_no_report_data() {
        let result = generate_tdx_report_data(None, "IXUKoBO1XEFBPwopN4sY".to_string());
        assert!(result.is_ok());
    }

    #[test]
    //generate_tdx_report allow empty report data string
    fn generate_tdx_report_data_report_data_size_0() {
        let result =
            generate_tdx_report_data(Some("".to_string()), "IXUKoBO1XEFBPwopN4sY".to_string());
        assert!(result.is_ok());
    }

    #[test]
    //generate_tdx_report allow 8 bytes report data string
    fn generate_tdx_report_data_report_data_size_8() {
        let result = generate_tdx_report_data(
            Some("YWJjZGVmZw==".to_string()),
            "IXUKoBO1XEFBPwopN4sY".to_string(),
        );
        assert!(result.is_ok());
    }

    #[test]
    //generate_tdx_report allow 48 bytes report data string
    fn generate_tdx_report_data_size_report_data_size_48() {
        // this one should be standard 48 bytes base64 encoded report data
        // "MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4" is base64 of "123456781234567812345678123456781234567812345678", 48 bytes
        let result = generate_tdx_report_data(
            Some("MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4".to_string()),
            "IXUKoBO1XEFBPwopN4sY".to_string(),
        );
        assert!(result.is_ok());
    }

    #[test]
    //generate_tdx_report require report data string is base64 encoded
    fn generate_tdx_report_data_report_data_not_base64_encoded() {
        //coming in report data should always be base64 encoded
        let result = generate_tdx_report_data(
            Some("XD^%*!x".to_string()),
            "IXUKoBO1XEFBPwopN4sY".to_string(),
        );
        assert!(result.is_err());
    }

    #[test]
    //generate_tdx_report require nonce string is base64 encoded
    fn generate_tdx_report_data_nonce_not_base64_encoded() {
        //coming in nonce should always be base64 encoded
        let result = generate_tdx_report_data(
            Some("IXUKoBO1XEFBPwopN4sY".to_string()),
            "XD^%*!x".to_string(),
        );
        assert!(result.is_err());
    }

    #[test]
    //generate_tdx_report require nonce string is base64 encoded
    fn generate_tdx_report_data_nonce_short_not_base64_encoded() {
        //coming in nonce should always be base64 encoded
        let result =
            generate_tdx_report_data(Some("IXUKoBO1XEFBPwopN4sY".to_string()), "123".to_string());
        assert!(result.is_err());
    }

    #[test]
    //generate_tdx_report require report data string is base64 encoded
    fn generate_tdx_report_data_report_data_short_not_base64_encoded() {
        //coming in report data should always be base64 encoded
        let result =
            generate_tdx_report_data(Some("123".to_string()), "IXUKoBO1XEFBPwopN4sY".to_string());
        assert!(result.is_err());
    }

    #[test]
    //generate_tdx_report check result as expected
    //original report_data = "abcdefgh", orginal nonce = "12345678"
    fn generate_tdx_report_data_report_data_nonce_base64_encoded_as_expected() {
        let result =
            generate_tdx_report_data(Some("YWJjZGVmZw==".to_string()), "MTIzNDU2Nzg=".to_string())
                .unwrap();
        let expected_hash = [
            93, 71, 28, 83, 115, 189, 166, 130, 87, 137, 126, 119, 140, 209, 163, 215, 13, 175,
            225, 101, 64, 195, 196, 202, 15, 37, 166, 241, 141, 49, 128, 157, 164, 132, 67, 50, 9,
            32, 162, 89, 243, 191, 177, 131, 4, 159, 156, 104, 11, 193, 18, 217, 92, 215, 194, 98,
            145, 191, 211, 85, 187, 118, 39, 80,
        ];
        let generated_hash = base64::decode(result).unwrap();
        assert_eq!(generated_hash, expected_hash);
    }

    #[test]
    //generate_tdx_report allow long report data string
    fn generate_tdx_report_data_long_tdx_report_data() {
        let result = generate_tdx_report_data(
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
    //generate_tdx_report allow long nonce string
    fn generate_tdx_report_data_long_nonce() {
        let result = generate_tdx_report_data(
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
    //generate_tdx_report_data generated report data is 64 bytes
    fn generate_tdx_report_data_report_data_is_64_bytes() {
        let report_data_hashed = match generate_tdx_report_data(
            Some("MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4".to_string()),
            "IXUKoBO1XEFBPwopN4sY".to_string(),
        ) {
            Ok(r) => r,
            Err(_) => todo!(),
        };
        let generated_hash_len = base64::decode(report_data_hashed).unwrap().len();
        assert_eq!(generated_hash_len, 64);
    }

    #[test]
    //TDX ENV required: tdx_get_quote allow empty nonce
    fn tdx_get_quote_empty_nonce() {
        let result = get_tdx_quote(Some("YWJjZGVmZw==".to_string()), "".to_string());
        assert!(result.is_ok());
    }

    #[test]
    //TDX ENV required: tdx_get_quote allow 0 bytes report data string
    fn tdx_get_quote_report_data_size_0() {
        let result = get_tdx_quote(Some("".to_string()), "IXUKoBO1XEFBPwopN4sY".to_string());
        assert!(result.is_ok());
    }

    #[test]
    //TDX ENV required: tdx_get_quote allow 8 bytes report data string
    fn tdx_get_quote_report_data_size_8() {
        // "YWJjZGVmZw==" is base64 of "abcdefg", 8 bytes
        let result = get_tdx_quote(
            Some("YWJjZGVmZw==".to_string()),
            "IXUKoBO1XEFBPwopN4sY".to_string(),
        );
        assert!(result.is_ok());
    }

    #[test]
    //TDX ENV required: tdx_get_quote allow 48 bytes report data string
    fn tdx_get_quote_report_data_size_48() {
        let result = get_tdx_quote(
            Some("MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4".to_string()),
            "IXUKoBO1XEFBPwopN4sY".to_string(),
        );
        assert!(result.is_ok());
    }

    #[test]
    //TDX ENV required: tdx_get_quote allow optional report data
    fn tdx_get_quote_report_data_null() {
        let result = get_tdx_quote(None, "IXUKoBO1XEFBPwopN4sY".to_string());
        assert!(result.is_ok());
    }

    #[test]
    //TDX ENV required: tdx_get_quote require report data string is base64 encoded
    fn tdx_get_quote_report_data_not_base64_encoded() {
        let result = get_tdx_quote(
            Some("XD^%*!x".to_string()),
            "IXUKoBO1XEFBPwopN4sY".to_string(),
        );
        assert!(result.is_err());
    }

    #[test]
    //TDX ENV required: tdx_get_quote require nonce string is base64 encoded
    fn tdx_get_quote_nonce_not_base64_encoded() {
        let result = get_tdx_quote(
            Some("IXUKoBO1XEFBPwopN4sY".to_string()),
            "XD^%*!x".to_string(),
        );
        assert!(result.is_err());
    }

    #[test]
    //TDX ENV required: tdx_get_quote allow long report data string
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
    //TDX ENV required: tdx_get_quote allow long nonce string
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
    //TDX ENV required: get_tdx_quote return non-empty encoded quote string
    fn tdx_get_quote_report_data_encoded_quote_is_not_0_bytes() {
        let quote = match get_tdx_quote(
            Some("MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4MTIzNDU2NzgxMjM0NTY3ODEyMzQ1Njc4".to_string()),
            "IXUKoBO1XEFBPwopN4sY".to_string(),
        ) {
            Ok(r) => r,
            Err(_) => todo!(),
        };
        assert_ne!(quote.len(), 0);
    }

    #[test]
    //get_quote does not allow tee type beyond TDX/SEV/TPM
    fn get_quote_wrong_tee_type() {
        let result = get_quote(
            TeeType::PLAIN,
            "".to_string(),
            "IXUKoBO1XEFBPwopN4sY".to_string(),
        );
        assert!(result.is_err());
    }

    #[test]
    //get_quote does not support SEV for now
    fn get_quote_sev_tee_type() {
        //does not allow tee type beyond TDX/SEV/TPM
        let result = get_quote(
            TeeType::SEV,
            "".to_string(),
            "IXUKoBO1XEFBPwopN4sY".to_string(),
        );
        assert!(result.is_err());
    }

    #[test]
    //get_quote does not support TPM for now
    fn get_quote_tpm_tee_type() {
        //does not allow tee type beyond TDX/SEV/TPM
        let result = get_quote(
            TeeType::TPM,
            "".to_string(),
            "IXUKoBO1XEFBPwopN4sY".to_string(),
        );
        assert!(result.is_err());
    }

    #[test]
    //get_quote support TDX now
    fn get_quote_tdx_tee_type() {
        //does not allow tee type beyond TDX/SEV/TPM
        let result = get_quote(
            TeeType::TDX,
            "".to_string(),
            "IXUKoBO1XEFBPwopN4sY".to_string(),
        );
        assert!(result.is_ok());
    }
}
