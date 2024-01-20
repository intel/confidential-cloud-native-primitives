
use crate::{
    ccnp_pb::{
        ccnp_server::Ccnp, 
        GetMeasurementRequest, GetMeasurementResponse, 
        Level, 
        GetEventlogRequest, GetEventlogResponse, 
        GetReportRequest, GetReportResponse,
        TcgEvent, TcgDigest
    }, handler::handle, tee::ITee
};
use std::result::Result::Ok;
use anyhow::Result;
use cctrusted_base::tcg::EventLogEntry;

use tonic::{Request, Response, Status};
pub struct Service {
    tee: Box<dyn ITee>
}

impl Service {
    pub fn new(tee: Box<dyn ITee>) -> Self {
        Service {
            tee
        }
    }
    pub fn tee_name(&self) -> String {
        self.tee.name()
    }
}

#[tonic::async_trait]
impl Ccnp for Service {
    async fn get_measurement(
        &self,
        request: Request<GetMeasurementRequest>,
    ) -> Result<Response<GetMeasurementResponse>, Status> {
        handle(&self.tee, request, 
            |tee: &Box<dyn ITee>, req: GetMeasurementRequest| -> Result<GetMeasurementResponse, Status> {
                let level = Level::from_i32(req.level)
                    .map_or_else(
                        || {
                            Err(Status::invalid_argument("Invalid Measurement Level"))
                        }, 
                        |val| {
                            Ok(val)
                        }
                    )?;
                tee.measurement(level, req.index as u8)
                    .map_or_else(
                        |e|{
                            Err(Status::internal(e.to_string()))
                        }, 
                        |val|{
                            Ok(
                                GetMeasurementResponse{
                                    measurement: val
                                }
                            )
                        }
                    )
            }
        )
    }
        
    async fn get_eventlog(
        &self,
        request: Request<GetEventlogRequest>,
    ) -> Result<Response<GetEventlogResponse>, Status> {
        handle(&self.tee, request, 
            |tee: &Box<dyn ITee>, req: GetEventlogRequest| -> Result<GetEventlogResponse, Status> {
                let level = Level::from_i32(req.level)
                    .map_or_else(
                        || {
                            Err(Status::invalid_argument("Invalid Eventlog Level"))
                        }, 
                        |val| {
                            Ok(val)
                        }
                    )?;
                tee.eventlog(level, req.start, req.count)
                    .map_or_else(
                        |e|{
                            Err(Status::internal(e.to_string()))
                        }, 
                        |entries|{
                            let mut eventlogs: Vec<TcgEvent> = vec![];
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
                                }
                            }
                            Ok(
                                GetEventlogResponse{
                                    events: eventlogs
                                } 
                            )
                        }
                    )
            }            
        )
        
    }

    async fn get_report(&self, request: Request<GetReportRequest>) -> Result<Response<GetReportResponse>, Status> {
        handle(&self.tee, request, 
            |tee: &Box<dyn ITee>, req: GetReportRequest| -> Result<GetReportResponse, Status> {
                let level = Level::from_i32(req.level)
                    .map_or_else(
                        || {
                            Err(Status::invalid_argument("Invalid Report Level"))
                        }, 
                        |val| {
                            Ok(val)
                        }
                    )?;
                tee.report(level, Some(req.user_data), Some(req.nonce))
                    .map_or_else(
                        |e|{
                            Err(Status::internal(e.to_string()))
                        }, 
                        |val|{
                            Ok(
                                GetReportResponse{
                                    report: val,
                                }
                            )
                        }
                    
                    )
            }
        )
    }
    
}

// todo: unit test