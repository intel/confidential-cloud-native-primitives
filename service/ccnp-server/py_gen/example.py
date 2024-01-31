from cli import CCNPClient
from cctrusted_vm import CCTrustedVmSdk
from report import Report

if __name__ == "__main__":  
    cli = CCNPClient()

    print("================ Report   Section { ================")  
    resp = cli.GetReport("","")
    print(resp.report)
    print("================ Report   Section } ================") 

    print("================ Measure  Section { ================") 
    alg = CCTrustedVmSdk.inst().get_default_algorithms()
    resp = cli.GetMeasurement(0, alg.alg_id)
    print(resp.measurement)
    print("================ Measure  Section } ================") 

    print("================ Eventlog Section { ================") 
    resp = cli.GetEventlog(1, 3)
    print(resp.events)
    print("================ Eventlog Section } ================") 