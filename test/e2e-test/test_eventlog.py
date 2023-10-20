import pytest
import base64
from ccnp import Eventlog
from ccnp import EventlogType
from pytdxattest.actor  import  TDEventLogActor
from pytdxattest.ccel  import  CCEL


class Testeventlog:
        def test_every_event(self):
                logs = Eventlog.get_platform_eventlog()
                ccelobj = CCEL.create_from_acpi_file()
                if ccelobj is None:
                        return
                ccelobj.dump()
                actor = TDEventLogActor(ccelobj.log_area_start_address,ccelobj.log_area_minimum_length)
                actor.process()
                for i in range(0,len(logs)):
                        assert base64.b64decode(logs[i].event) == actor._event_logs[i]._event

        def test_every_digest(self):
                logs = Eventlog.get_platform_eventlog()
                ccelobj = CCEL.create_from_acpi_file()
                if ccelobj is None:
                        return
                ccelobj.dump()
                actor = TDEventLogActor(ccelobj.log_area_start_address,ccelobj.log_area_minimum_length)
                actor.process()
                assert len(logs)==len(actor._event_logs)
                for i in range(0,len(logs)):
                        digest=[]
                        for j in range(0,len((actor._event_logs[i]._digests)[0])):
                                digest.append((actor._event_logs[i]._digests)[0][j])
                        assert str(digest).replace(",","") == logs[i].digest

