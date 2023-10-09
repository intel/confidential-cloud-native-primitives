/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
*/

#![allow(non_camel_case_types)]

use anyhow::*;
use nix::*;
use std::convert::TryInto;
use std::fs::File;
use std::mem;
use std::os::unix::io::AsRawFd;
use std::path::Path;
use std::ptr;
use std::result::Result;
use std::result::Result::Ok;

#[repr(C)]
pub struct tdx_1_0_report_req {
    subtype: u8,     // Subtype of TDREPORT: fixed as 0 by TDX Module specification
    reportdata: u64, // User-defined REPORTDATA to be included into TDREPORT
    rpd_len: u32,    // Length of the REPORTDATA: fixed as 64 bytes by the TDX Module specification
    tdreport: u64,   // TDREPORT output from TDCALL[TDG.MR.REPORT]
    tdr_len: u32,    // Length of the TDREPORT: fixed as 1024 bytes by the TDX Module specification
}

#[repr(C)]
pub struct tdx_1_5_report_req {
    reportdata: [u8; REPORT_DATA_LEN as usize], // User buffer with REPORTDATA to be included into TDREPORT
    tdreport: [u8; TDX_REPORT_LEN as usize], // User buffer to store TDREPORT output from TDCALL[TDG.MR.REPORT]
}

#[repr(C)]
pub struct qgs_msg_header {
    major_version: u16, // TDX major version
    minor_version: u16, // TDX minor version
    msg_type: u32,      // GET_QUOTE_REQ or GET_QUOTE_RESP
    size: u32,          // size of the whole message, include this header, in byte
    error_code: u32,    // used in response only
}

#[repr(C)]
pub struct qgs_msg_get_quote_req {
    header: qgs_msg_header,                        // header.type = GET_QUOTE_REQ
    report_size: u32,                              // cannot be 0
    id_list_size: u32,                             // length of id_list, in byte, can be 0
    report_id_list: [u8; TDX_REPORT_LEN as usize], // report followed by id list
}

#[repr(C)]
pub struct tdx_quote_hdr {
    version: u64,                       // Quote version, filled by TD
    status: u64,                        // Status code of Quote request, filled by VMM
    in_len: u32,                        // Length of TDREPORT, filled by TD
    out_len: u32,                       // Length of Quote, filled by VMM
    data_len_be_bytes: [u8; 4],         // big-endian 4 bytes indicate the size of data following
    data: [u8; TDX_QUOTE_LEN as usize], // Actual Quote data or TDREPORT on input
}

#[repr(C)]
pub struct tdx_quote_req {
    buf: u64, // Pass user data that includes TDREPORT as input. Upon successful completion of IOCTL, output is copied back to the same buffer
    len: u64, // Length of the Quote buffer
}

#[repr(C)]
pub struct qgs_msg_get_quote_resp {
    header: qgs_msg_header,        // header.type = GET_QUOTE_RESP
    selected_id_size: u32,         // can be 0 in case only one id is sent in request
    quote_size: u32,               // length of quote_data, in byte
    id_quote: [u8; TDX_QUOTE_LEN], // selected id followed by quote
}

pub enum TdxVersion {
    TDX_1_0,
    TDX_1_5,
}

pub enum TdxOperation {
    TDX_GET_TD_REPORT = 1,
    TDX_1_0_GET_QUOTE = 2,
    TDX_1_5_GET_QUOTE = 4,
}

const REPORT_DATA_LEN: u32 = 64;
const TDX_REPORT_LEN: u32 = 1024;
const TDX_QUOTE_LEN: usize = 4 * 4096;

pub struct TdxInfo {
    tdx_version: TdxVersion,
    device_node: File,
}

impl TdxInfo {
    fn new(_tdx_version: TdxVersion, _device_node: File) -> Self {
        TdxInfo {
            tdx_version: _tdx_version,
            device_node: _device_node,
        }
    }
}

fn get_tdx_version() -> TdxVersion {
    if Path::new("/dev/tdx-guest").exists() {
        TdxVersion::TDX_1_0
    } else if Path::new("/dev/tdx_guest").exists() {
        TdxVersion::TDX_1_5
    } else if Path::new("/dev/tdx-attest").exists() {
        panic!("get_tdx_version: Deprecated device node /dev/tdx-attest, please upgrade to use /dev/tdx-guest or /dev/tdx_guest");
    } else {
        panic!("get_tdx_version: no TDX device found!");
    }
}

pub fn get_td_report(report_data: String) -> Result<Vec<u8>, anyhow::Error> {
    //detect TDX version
    let tdx_info = match get_tdx_version() {
        TdxVersion::TDX_1_0 => {
            let device_node = match File::options()
                .read(true)
                .write(true)
                .open("/dev/tdx-guest")
            {
                Err(e) => {
                    return Err(anyhow!(
                        "[get_td_report] Fail to open {}: {:?}",
                        "/dev/tdx-guest",
                        e
                    ))
                }
                Ok(fd) => fd,
            };
            TdxInfo::new(TdxVersion::TDX_1_0, device_node)
        }
        TdxVersion::TDX_1_5 => {
            let device_node = match File::options()
                .read(true)
                .write(true)
                .open("/dev/tdx_guest")
            {
                Err(e) => {
                    return Err(anyhow!(
                        "[get_td_report] Fail to open {}: {:?}",
                        "/dev/tdx_guest",
                        e
                    ))
                }
                Ok(fd) => fd,
            };
            TdxInfo::new(TdxVersion::TDX_1_5, device_node)
        }
    };

    match tdx_info.tdx_version {
        TdxVersion::TDX_1_0 => match get_tdx_1_0_report(tdx_info.device_node, report_data) {
            Err(e) => return Err(anyhow!("[get_td_report] Fail to get TDX report: {:?}", e)),
            Ok(report) => Ok(report),
        },
        TdxVersion::TDX_1_5 => match get_tdx_1_5_report(tdx_info.device_node, report_data) {
            Err(e) => return Err(anyhow!("[get_td_report] Fail to get TDX report: {:?}", e)),
            Ok(report) => Ok(report),
        },
    }
}

fn get_tdx_1_0_report(device_node: File, report_data: String) -> Result<Vec<u8>, anyhow::Error> {
    let report_data_bytes = match base64::decode(report_data) {
        Ok(v) => v,
        Err(e) => return Err(anyhow!("report data is not base64 encoded: {:?}", e)),
    };

    //prepare get TDX report request data
    let mut report_data_array: [u8; REPORT_DATA_LEN as usize] = [0; REPORT_DATA_LEN as usize];
    report_data_array.copy_from_slice(&report_data_bytes[0..]);
    let td_report: [u8; TDX_REPORT_LEN as usize] = [0; TDX_REPORT_LEN as usize];

    //build the request
    let request = tdx_1_0_report_req {
        subtype: 0 as u8,
        reportdata: ptr::addr_of!(report_data_array) as u64,
        rpd_len: REPORT_DATA_LEN,
        tdreport: ptr::addr_of!(td_report) as u64,
        tdr_len: TDX_REPORT_LEN,
    };

    //build the operator code
    ioctl_readwrite!(
        get_report_1_0_ioctl,
        b'T',
        TdxOperation::TDX_GET_TD_REPORT,
        u64
    );

    //apply the ioctl command
    match unsafe {
        get_report_1_0_ioctl(device_node.as_raw_fd(), ptr::addr_of!(request) as *mut u64)
    } {
        Err(e) => {
            return Err(anyhow!(
                "[get_tdx_1_0_report] Fail to get TDX report: {:?}",
                e
            ))
        }
        Ok(_) => (),
    };

    Ok(td_report.to_vec())
}

fn get_tdx_1_5_report(device_node: File, report_data: String) -> Result<Vec<u8>, anyhow::Error> {
    let report_data_bytes = match base64::decode(report_data) {
        Ok(v) => v,
        Err(e) => return Err(anyhow!("report data is not base64 encoded: {:?}", e)),
    };

    //prepare get TDX report request data
    let mut request = tdx_1_5_report_req {
        reportdata: [0; REPORT_DATA_LEN as usize],
        tdreport: [0; TDX_REPORT_LEN as usize],
    };
    request.reportdata.copy_from_slice(&report_data_bytes[0..]);

    //build the operator code
    ioctl_readwrite!(
        get_report_1_5_ioctl,
        b'T',
        TdxOperation::TDX_GET_TD_REPORT,
        tdx_1_5_report_req
    );

    //apply the ioctl command
    match unsafe {
        get_report_1_5_ioctl(
            device_node.as_raw_fd(),
            ptr::addr_of!(request) as *mut tdx_1_5_report_req,
        )
    } {
        Err(e) => {
            return Err(anyhow!(
                "[get_tdx_1_5_report] Fail to get TDX report: {:?}",
                e
            ))
        }
        Ok(_) => (),
    };

    Ok(request.tdreport.to_vec())
}

fn generate_qgs_quote_msg(report: [u8; TDX_REPORT_LEN as usize]) -> qgs_msg_get_quote_req {
    //build quote service message header to be used by QGS
    let qgs_header = qgs_msg_header {
        major_version: 1,
        minor_version: 0,
        msg_type: 0,
        size: 16 + 8 + TDX_REPORT_LEN, // header + report_size and id_list_size + TDX_REPORT_LEN
        error_code: 0,
    };

    //build quote service message body to be used by QGS
    let mut qgs_request = qgs_msg_get_quote_req {
        header: qgs_header,
        report_size: TDX_REPORT_LEN,
        id_list_size: 0,
        report_id_list: [0; TDX_REPORT_LEN as usize],
    };

    qgs_request.report_id_list.copy_from_slice(&report[0..]);

    qgs_request
}

pub fn get_tdx_quote(report_data: String) -> Result<Vec<u8>, anyhow::Error> {
    //retrive TDX report
    let report_data_vec = match get_td_report(report_data) {
        Err(e) => return Err(anyhow!("[get_tdx_quote] Fail to get TDX report: {:?}", e)),
        Ok(report) => report,
    };
    let report_data_array: [u8; TDX_REPORT_LEN as usize] = match report_data_vec.try_into() {
        Ok(r) => r,
        Err(e) => return Err(anyhow!("[get_tdx_quote] Wrong TDX report format: {:?}", e)),
    };

    //build QGS request message
    let qgs_msg = generate_qgs_quote_msg(report_data_array);

    let tdx_info = match get_tdx_version() {
        TdxVersion::TDX_1_0 => {
            let device_node = match File::options()
                .read(true)
                .write(true)
                .open("/dev/tdx-guest")
            {
                Err(e) => {
                    return Err(anyhow!(
                        "[get_tdx_quote] Fail to open {}: {:?}",
                        "/dev/tdx-guest",
                        e
                    ))
                }
                Ok(fd) => fd,
            };
            TdxInfo::new(TdxVersion::TDX_1_0, device_node)
        }
        TdxVersion::TDX_1_5 => {
            let device_node = match File::options()
                .read(true)
                .write(true)
                .open("/dev/tdx_guest")
            {
                Err(e) => {
                    return Err(anyhow!(
                        "[get_tdx_quote] Fail to open {}: {:?}",
                        "/dev/tdx_guest",
                        e
                    ))
                }
                Ok(fd) => fd,
            };
            TdxInfo::new(TdxVersion::TDX_1_5, device_node)
        }
    };

    //build quote generation request header
    let mut quote_header = tdx_quote_hdr {
        version: 1,
        status: 0,
        in_len: (mem::size_of_val(&qgs_msg) + 4) as u32,
        out_len: 0,
        data_len_be_bytes: (1048 as u32).to_be_bytes(),
        data: [0; TDX_QUOTE_LEN as usize],
    };

    let qgs_msg_bytes = unsafe {
        let ptr = &qgs_msg as *const qgs_msg_get_quote_req as *const u8;
        std::slice::from_raw_parts(ptr, mem::size_of::<qgs_msg_get_quote_req>())
    };
    quote_header.data[0..(16 + 8 + TDX_REPORT_LEN) as usize]
        .copy_from_slice(&qgs_msg_bytes[0..((16 + 8 + TDX_REPORT_LEN) as usize)]);

    let request = tdx_quote_req {
        buf: ptr::addr_of!(quote_header) as u64,
        len: TDX_QUOTE_LEN as u64,
    };

    //build the operator code and apply the ioctl command
    match tdx_info.tdx_version {
        TdxVersion::TDX_1_0 => {
            ioctl_read!(
                get_quote_1_0_ioctl,
                b'T',
                TdxOperation::TDX_1_0_GET_QUOTE,
                u64
            );
            match unsafe {
                get_quote_1_0_ioctl(
                    tdx_info.device_node.as_raw_fd(),
                    ptr::addr_of!(request) as *mut u64,
                )
            } {
                Err(e) => return Err(anyhow!("[get_tdx_quote] Fail to get TDX quote: {:?}", e)),
                Ok(_r) => _r,
            };
        }
        TdxVersion::TDX_1_5 => {
            ioctl_read!(
                get_quote_1_5_ioctl,
                b'T',
                TdxOperation::TDX_1_5_GET_QUOTE,
                tdx_quote_req
            );
            match unsafe {
                get_quote_1_5_ioctl(
                    tdx_info.device_node.as_raw_fd(),
                    ptr::addr_of!(request) as *mut tdx_quote_req,
                )
            } {
                Err(e) => return Err(anyhow!("[get_tdx_quote] Fail to get TDX quote: {:?}", e)),
                Ok(_r) => _r,
            };
        }
    };

    //inspect the response and retrive quote data
    let out_len = quote_header.out_len;
    let qgs_msg_resp_size =
        unsafe { std::mem::transmute::<[u8; 4], u32>(quote_header.data_len_be_bytes) }.to_be();

    let qgs_msg_resp = unsafe {
        let raw_ptr = ptr::addr_of!(quote_header.data) as *mut qgs_msg_get_quote_resp;
        raw_ptr.as_mut().unwrap() as &mut qgs_msg_get_quote_resp
    };

    if out_len - qgs_msg_resp_size != 4 {
        return Err(anyhow!(
            "[get_tdx_quote] Fail to get TDX quote: wrong TDX quote size!"
        ));
    }

    if qgs_msg_resp.header.major_version != 1
        || qgs_msg_resp.header.minor_version != 0
        || qgs_msg_resp.header.msg_type != 1
        || qgs_msg_resp.header.error_code != 0
    {
        return Err(anyhow!(
            "[get_tdx_quote] Fail to get TDX quote: QGS response error!"
        ));
    }

    Ok(qgs_msg_resp.id_quote[0..(qgs_msg_resp.quote_size as usize)].to_vec())
}

#[cfg(test)]
mod tdx_attest_tests {
    use super::*;

    #[test]
    //TDX ENV required: call get_td_report and verify report data embedded in quote
    fn get_td_report_verify_report_data() {
        let report_data = "XUccU3O9poJXiX53jNGj1w2v4WVAw8TKDyWm8Y0xgJ2khEMyCSCiWfO/sYMEn5xoC8ES2VzXwmKRv9NVu3YnUA==";
        let report = get_td_report(report_data.to_string()).unwrap();

        let expected_report_data = [
            93, 71, 28, 83, 115, 189, 166, 130, 87, 137, 126, 119, 140, 209, 163, 215, 13, 175,
            225, 101, 64, 195, 196, 202, 15, 37, 166, 241, 141, 49, 128, 157, 164, 132, 67, 50, 9,
            32, 162, 89, 243, 191, 177, 131, 4, 159, 156, 104, 11, 193, 18, 217, 92, 215, 194, 98,
            145, 191, 211, 85, 187, 118, 39, 80,
        ];

        let mut report_data_in_report: [u8; 64 as usize] = [0; 64 as usize];
        report_data_in_report.copy_from_slice(&report[128..192]);
        assert_eq!(report_data_in_report, expected_report_data);
    }

    #[test]
    //TDX ENV required: call tdx_get_quote and verify report data embedded in quote
    fn get_tdx_quote_verify_report_data() {
        let report_data = "XUccU3O9poJXiX53jNGj1w2v4WVAw8TKDyWm8Y0xgJ2khEMyCSCiWfO/sYMEn5xoC8ES2VzXwmKRv9NVu3YnUA==";
        let quote = get_tdx_quote(report_data.to_string()).unwrap();

        let expected_report_data = [
            93, 71, 28, 83, 115, 189, 166, 130, 87, 137, 126, 119, 140, 209, 163, 215, 13, 175,
            225, 101, 64, 195, 196, 202, 15, 37, 166, 241, 141, 49, 128, 157, 164, 132, 67, 50, 9,
            32, 162, 89, 243, 191, 177, 131, 4, 159, 156, 104, 11, 193, 18, 217, 92, 215, 194, 98,
            145, 191, 211, 85, 187, 118, 39, 80,
        ];

        let mut report_data_in_quote: [u8; 64 as usize] = [0; 64 as usize];
        report_data_in_quote.copy_from_slice(&quote[568..632]);
        assert_eq!(report_data_in_quote, expected_report_data);
    }
}
