use anyhow::Result;

use cctrusted_vm::cvm::get_cvm_type;
use cctrusted_base::{cc_type::{self, CcType}, tcg::EventLogEntry};

use crate::ccnp_pb::Level;

pub mod tdx;

pub trait ITee: Sync + Send { 
    fn init(&mut self, tee_type: CcType);
    fn name(&self) -> String;

    fn eventlog(&self, 
        level: Level, 
        start: u32, 
        count: u32
    ) -> Result<Vec<EventLogEntry>>;
    fn measurement(&self, 
        level: Level,
        index: u8
    ) -> Result<Vec<u8>>;
    fn report(&self, 
        level: Level, 
        user_data: Option<String>, 
        nonce: Option<String>
    ) -> Result<Vec<u8>>;
    
}


pub fn get_available_tee() -> Option<Box<dyn ITee>> {
    let cvm_type = get_cvm_type();
    return match cvm_type.tee_type {
        cc_type::TeeType::TDX  => {
            let mut tdx_tee= tdx::TDX::default();
            tdx_tee.init(cvm_type);
            Some(
                Box::new(tdx_tee)
            )
        }
        _ => None
    }
}