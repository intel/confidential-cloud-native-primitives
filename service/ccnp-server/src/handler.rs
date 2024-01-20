use std::fmt::Debug;
use anyhow::Result;
use log::info;
use tonic::{Request, Response, Status};
use uuid::Uuid;

use crate::tee::ITee;

pub fn handle<ReqType, RespType, F>
    (tee: &Box<dyn ITee>, request: Request<ReqType>, operate: F)
    -> Result<Response<RespType>, Status>
    where ReqType: Debug, RespType: Debug, 
        F: Fn(&Box<dyn ITee>, ReqType) -> Result<RespType, Status>
{
    let req: ReqType = request.into_inner();  
    let id = Uuid::new_v4();
    info!("Request  IN ---> id:{:}: {:?}", id, req);

    let res = operate(tee, req);
    return match res {
        Ok(resp) => {
            info!("Response OK <--- id:{:}", id);
            Ok(Response::new(resp))
        }
        Err(e) => {
            info!("Response ER <--- id:{:}: {:?}", id, e.to_string());
            Err(e)
        }
    }
    
}