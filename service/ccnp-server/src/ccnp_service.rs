use crate::ccnp_pb::{
        ccnp_server::Ccnp, 
        GetMeasurementRequest, GetMeasurementResponse,  
        GetEventlogRequest, GetEventlogResponse, 
        GetReportRequest, GetReportResponse,
        TcgEvent, TcgDigest
    };
use std::result::Result::Ok;
use anyhow::Result;
use cctrusted_base::{api::CCTrustedApi, api_data::ExtraArgs, tcg::EventLogEntry};

use cctrusted_vm::sdk::API;
use tonic::{Request, Response, Status};

#[derive(Default)]
pub struct Service {}

#[tonic::async_trait]
impl Ccnp for Service {
    async fn get_measurement(
        &self,
        request: Request<GetMeasurementRequest>,
    ) -> Result<Response<GetMeasurementResponse>, Status> {
        let req = request.into_inner();
        API::get_cc_measurement(req.index as u8, req.algo_id as u16)
            .map_or_else(
                |e|{
                    Err(Status::internal(e.to_string()))
                }, 
                |val|{
                    Ok(
                        Response::new(
                            GetMeasurementResponse{
                                measurement: val.get_hash()
                            }
                        )
                    )
                }
            )
    }
        
    async fn get_eventlog(
        &self,
        request: Request<GetEventlogRequest>,
    ) -> Result<Response<GetEventlogResponse>, Status> {
        let req = request.into_inner();
        let mut eventlogs: Vec<TcgEvent> = vec![];
        let entries = match API::get_cc_eventlog(Some(req.start), Some(req.count)) {
            Ok(val) => Ok(val),
            Err(e)  => Err(Status::internal(e.to_string())),
        }?;

        for entry in entries {
            match entry {
                EventLogEntry::TcgImrEvent(event) => {
                    let mut digests: Vec<TcgDigest> = vec![]; 
                    for d in event.digests {
                        digests.push(TcgDigest{
                            algo_id: d.algo_id as u32,
                            hash: d.hash,
                        })
                    }
                    eventlogs.push(TcgEvent{
                        imr_index: event.imr_index,
                        event_type: event.event_type,
                        event_size: event.event_size,
                        event: event.event,
                        digest: vec![], 
                        digests, 
                    })
                }
                EventLogEntry::TcgPcClientImrEvent(event) => {
                    eventlogs.push(TcgEvent{
                        imr_index: event.imr_index,
                        event_type: event.event_type,
                        event_size: event.event_size,
                        event: event.event,
                        digest: event.digest.to_vec(), 
                        digests: vec![], 
                    })
                }
                EventLogEntry::TcgCanonicalEvent(_) => {
                    return Err(Status::internal("Unsupported Event: Tcg Canonical Event"))
                },
            }
        }

        Ok(Response::new(GetEventlogResponse{events: eventlogs}))
        
    }

    async fn get_report(&self, request: Request<GetReportRequest>) -> Result<Response<GetReportResponse>, Status> {
        let req = request.into_inner();
        API::get_cc_report(Some(req.nonce), Some(req.user_data), ExtraArgs {})
            .map_or_else(
                |e|{
                    Err(Status::internal(e.to_string()))
                }, 
                |val|{
                    Ok(
                        Response::new(
                            GetReportResponse{
                                report: val.cc_report
                            }
                        )
                    )
                }
            )  
    }
    
}

// todo: unit test