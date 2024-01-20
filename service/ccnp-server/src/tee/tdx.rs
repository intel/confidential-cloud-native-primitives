use anyhow::{Result, Ok};
use cctrusted_base::{
    api::CCTrustedApi, 
    api_data::ExtraArgs, 
    tcg::EventLogEntry,
    cc_type::CcType
};
use cctrusted_vm::sdk::API;

use crate::ccnp_pb::Level;

use super::ITee;


pub struct TDX {
    tee_type: Option<CcType>,
}

impl Default for TDX {
    fn default() -> Self {
        Self { tee_type: None }
    }
}

impl ITee for TDX {
    fn init(&mut self, tee_type: CcType) {
        self.tee_type = Some(tee_type)
    }

    fn name(&self) -> String {
       self.tee_type.as_ref().map_or_else(
            || {
                "".to_string()
            }, 
            |val| {
                val.tee_type_str.clone()
            }
        ) 
    }

    fn eventlog(&self, 
        _level: Level, 
        start: u32, 
        count: u32
    ) -> Result<Vec<EventLogEntry>> {
        API::get_cc_eventlog(Some(start), Some(count))
    }

    fn measurement(&self, 
        _level: Level,
        index: u8
    ) -> Result<Vec<u8>> {
        let defalt_algo = API::get_default_algorithm()?;
        API::get_cc_measurement(index, defalt_algo.algo_id)
            .map_or_else(
                |e|{
                    Err(e)
                }, 
                |val|{
                    Ok(val.get_hash())
                }
            )
    }

    fn report(&self, 
        _level: Level, 
        user_data: Option<String>, 
        nonce: Option<String>
    ) -> Result<Vec<u8>> {
        API::get_cc_report(nonce, user_data, ExtraArgs {})
            .map_or_else(
                |e|{
                    Err(e)
                }, 
                |val|{
                    Ok(val.cc_report)
                }
            )
    }
}
