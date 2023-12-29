"""
This is a sample python script to do the following things:
1. Fetch boot time event logs using CCNP api
2. Fetch runtime event logs(from IMA) in kernel memory
3. Replay all event logs and re-calcuate the overall digest
4. Fetch IMR measurements using CCNP api
5. Compare the values and returns result

It also provides an option to verify selected runtime event logs.
Use "--verify-register-index" to specify the rtmr indexes to be verified. 
Use "-f" to specify the path of the reference event log file containing 
the event log with certain format, which is same as the IMA one:
register_index | template_hash | template | file_content
"""

import base64
import os
import logging
import argparse
import string

from hashlib import sha384
from typing import Dict
from ccnp import Measurement
from ccnp import MeasurementType
from ccnp import Eventlog

LOG = logging.getLogger(__name__)

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
    handlers=[
        logging.StreamHandler()
    ]
)

class BinaryBlob:
    """
    Manage the binary blob.
    """

    def __init__(self, data, base=0):
        self._data = data
        self._base_address = base

    @property
    def length(self):
        """Length of binary in bytes"""
        return len(self._data)

    @property
    def data(self):
        """Raw data of binary blob"""
        return self._data

    def to_hex_string(self):
        """To hex string"""
        return "".join(f"{b:02x}" % b for b in self._data)

    def get_bytes(self, pos, count):
        """Get bytes"""
        if count == 0:
            return None
        assert pos + count <= self.length
        return (self.data[pos:pos + count], pos + count)

    def dump(self):
        """Dump Hex value."""
        index = 0
        linestr = ""
        printstr = ""

        while index < self.length:
            if (index % 16) == 0:
                if len(linestr) != 0:
                    LOG.info("%s %s", linestr, printstr)
                    printstr = ''
                # line prefix string
                # pylint: disable=consider-using-f-string
                linestr = "{0:08X}  ".format(int(index / 16) * 16 + \
                    self._base_address)

            # pylint: disable=consider-using-f-string
            linestr += "{0:02X} ".format(self._data[index])
            if chr(self._data[index]) in set(string.printable) and \
               self._data[index] not in [0xC, 0xB, 0xA, 0xD, 0x9]:
                printstr += chr(self._data[index])
            else:
                printstr += '.'

            index += 1

        if (index % 16) != 0:
            blank = ""
            for _ in range(16 - index % 16):
                blank = blank + "   "
            LOG.info("%s%s %s", linestr, blank, printstr)
        elif index == self.length:
            LOG.info("%s %s", linestr, printstr)

class RTMR(BinaryBlob):
    """
    Data structure for RTMR registers.
    A RTMR register manages a 48-bytes (384-bits) hash value.
    """
    RTMR_COUNT = 4
    RTMR_LENGTH_BY_BYTES = 48

    def __init__(self, data: bytearray = bytearray(RTMR_LENGTH_BY_BYTES),
        base_addr=0):
        super().__init__(data, base_addr)

    def __eq__(self, other):
        bytearray_1, _ = self.get_bytes(0, RTMR.RTMR_LENGTH_BY_BYTES)
        bytearray_2, _ = other.get_bytes(0, RTMR.RTMR_LENGTH_BY_BYTES)

        return bytearray(bytearray_1) == bytearray(bytearray_2)

# pylint: disable=too-few-public-methods
class VerifyActor:
    """Actor to verify the RTMR
    """

    def _verify_single_rtmr(self, rtmr_index: int, rtmr_value_1: RTMR,
        rtmr_value_2: RTMR) -> bool:
        """Verify single RTMR value"""

        if rtmr_value_1 == rtmr_value_2:
            LOG.info("RTMR[%d] passed the verification.", rtmr_index)
            return True
        else:
            LOG.error("RTMR[%d] did not pass the verification", rtmr_index)
            return False

    def verify_rtmr(self, index_list: list, event_log_file: str,
                    ima_flag: bool) -> None:
        """Fetch RTMR measurement and event logs using CCNP API 
           and replay event log to do verification.
        """

        # 0. Print verify scope
        LOG.info("Step 0: List verify scope")
        LOG.info("Verifying RTMRs: [%s]\n", ','.join(str(x) for x in index_list))

        # 1. Check if IMA measurement event log exist at
        # /sys/kernel/security/integrity/ima/ascii_runtime_measurements
        LOG.info("Step 1: Check if IMA event logs exist in the system.")
        ima_measurement_file = \
                    "/run/security/integrity/ima/ascii_runtime_measurements"
        assert os.path.exists(ima_measurement_file), \
            f"Could not find the IMA measurement file {ima_measurement_file}"
        LOG.info("IMA event logs found in the system.\n")

        # 2. Init CCEventlogActor and collect event log
        #    and replay the RTMR value according to event log
        # pylint: disable-next=line-too-long
        LOG.info("Step 2: Collect boot time and runtime event logs and replay results.")
        cc_event_log_actor = CCEventLogActor()
        cc_event_log_actor.replay(index_list, ima_flag)

        # 3. Collect RTMR measurements using CCNP
        LOG.info("Step 3: Fetching measurements in RTMR.")
        rtmrs = []
        for index in index_list:
            LOG.info("==> Fetching measurements in RTMR[%d]", index)
            val = Measurement.get_platform_measurement(
                MeasurementType.TYPE_TDX_RTMR, None, index)
            rtmrs.append(val)
            LOG.info("RTMR[%d]: %s\n", index, base64.b64decode(val).hex())

        # 4. Verify individual RTMR value from CCNP fetching and recalculated from event log
        # pylint: disable-next=line-too-long
        LOG.info("Step 4: Verify individual RTMR value and re-calculated value from event logs")
        flag = True
        start = 0
        for index in index_list:
            flag = flag and self._verify_single_rtmr(
                index,
                cc_event_log_actor.get_rtmr_by_index(index),
                RTMR(bytearray(base64.b64decode(rtmrs[start]))))
            start += 1

        if flag:
            LOG.info("RTMR verify success.\n")
        else:
            LOG.info("RTMR verify failed. Skip selected event log verification if requested.\n")

        # 5. Verify selected digest according to file input
        if event_log_file is not None and flag:
            LOG.info("Step 5: Verify selected measurements from event logs.")
            cc_event_log_actor.verify_selected_runtime_measurement(event_log_file)

class CCEventLogActor:
    """Event log actor

    The actor to process event logs and do replay
    """

    RUNTIME_REGISTER = 2

    def __init__(self):
        self._boot_time_event_logs = []
        self._runtime_event_logs = []
        self._rtmrs:list[RTMR] = {}

    def _fetch_boot_time_event_logs(self):
        """Fetch cvm boot time event log using CCNP API.
        """
        LOG.info("==> Fetching boot time event logs using CCNP API")
        self._boot_time_event_logs = Eventlog.get_platform_eventlog()

    def _fetch_runtime_event_logs(self):
        """Fetch cvm runtime event log from IMA.
        """
        LOG.info("==> Fetching runtime event logs from IMA")
        ima_measurement_file = "/run/security/integrity/ima/ascii_runtime_measurements"
        with open(ima_measurement_file, encoding="utf-8") as f:
            num = 0
            for line in f:
                self._runtime_event_logs.append(line)
                elements = line.split(" ")
                self._dump_runtime_event_log(
                    int(elements[1]), 0x14,
                    "EV_IMA_NODE_MEASUREMENT_EVENT", elements[2])
                num = num + 1
        LOG.info("")

    @staticmethod
    def _replay_single_boot_time_rtmr(event_logs) -> RTMR:
        """Replay single RTMR for boot time events"""
        rtmr = bytearray(RTMR.RTMR_LENGTH_BY_BYTES)

        for event_log in event_logs:
            digest = list(map(int, event_log.digest.strip('[]').split(' ')))
            # pylint: disable-next=consider-using-f-string
            digest_hex = ''.join('{:02x}'.format(i) for i in digest)
            sha384_algo = sha384()
            sha384_algo.update(bytes.fromhex(rtmr.hex() + digest_hex))
            rtmr = sha384_algo.digest()

        return RTMR(rtmr)

    @staticmethod
    def _replay_runtime_rtmr(event_logs, base: RTMR) -> RTMR:
        """Replay runtime measurements based on the runtime event logs"""
        rtmr = bytearray(RTMR.RTMR_LENGTH_BY_BYTES)

        val = base.data.hex()
        for event_log in event_logs:
            elements = event_log.split(" ")
            extend_val = val + elements[2]
            sha384_algo = sha384()
            sha384_algo.update(bytes.fromhex(extend_val))
            val = sha384_algo.hexdigest()

        rtmr = sha384_algo.digest()
        return RTMR(rtmr)

    def get_rtmr_by_index(self, index: int) -> RTMR:
        """Get RTMR by TD register index"""
        return self._rtmrs[index]

    def _dump_boot_time_event_log(self, index, evt_type, type_str, digest):
        """Dump event log content"""
        LOG.info("--------------------Event Log Entry----------------------")
        LOG.info("RTMR index:         %d", index)
        LOG.info("Event type:         %d(%s)", evt_type, type_str)
        LOG.info("Digest:")
        digest_blob = BinaryBlob(digest)
        digest_blob.dump()


    def _dump_runtime_event_log(self, index, evt_type, type_str, digest):
        """Dump runtime event log content"""
        LOG.info("--------------------Event Log Entry----------------------")
        LOG.info("RTMR index:         %d", index)
        LOG.info("Event type:         %d(%s)", evt_type, type_str)
        LOG.info("Digest:")
        LOG.info("%s",digest)

    def replay(self, index_list:list, ima_flag:bool) -> Dict[int, RTMR]:
        """Replay event logs including boot time event logs and runtime event logs to
        generate RTMR values for verification.
        """
        self._fetch_boot_time_event_logs()

        boot_time_event_logs_by_index = {}
        for index in index_list:
            boot_time_event_logs_by_index[index] = []

        for event_log in self._boot_time_event_logs:
            if event_log.reg_idx in index_list:
                self._dump_boot_time_event_log(event_log.reg_idx,
                                               event_log.evt_type,
                                               event_log.evt_type_str,
                                               list(
                                                map(int, event_log.digest.strip('[]').split(' '))))
                LOG.info("")
                boot_time_event_logs_by_index[event_log.reg_idx].append(event_log)

        # replay boot time event logs and save replay results to dict
        rtmr_by_index = {}
        for rtmr_index, event_logs in boot_time_event_logs_by_index.items():
            rtmr_value = CCEventLogActor._replay_single_boot_time_rtmr(event_logs)
            rtmr_by_index[rtmr_index] = rtmr_value

        # runtime measurements are now extended into RTMR[2]
        # fetch and replay the runtime event logs if RTMR[2] included
        if CCEventLogActor.RUNTIME_REGISTER in index_list and ima_flag:
            self._fetch_runtime_event_logs()
            concat_rtmr_value = CCEventLogActor._replay_runtime_rtmr(
                self._runtime_event_logs, rtmr_by_index[2])
            rtmr_by_index[2] = concat_rtmr_value

        self._rtmrs = rtmr_by_index

    def verify_selected_runtime_measurement(self, digest_file: str):
        """Verify specific runtime event log entries against
        user provided hash golden values.
        """

        # open the digest file containing golden values
        # and find if the values exist in the event logs fetched.
        with open(digest_file, encoding="utf-8") as f:
            for line in f:
                flag = False
                ref_values = line.split(" ")
                LOG.info("==> Verifying digest: %s %s %s",
                         ref_values[1], ref_values[2], ref_values[4])
                for event_log in self._runtime_event_logs:
                    elements = event_log.split(" ")
                    if elements[2] == ref_values[1]:
                        flag = True
                        break
                if flag:
                    LOG.info("Verify success.")
                else:
                    LOG.info("Verify failed.")


if __name__ == "__main__":
    # pylint: disable-next=line-too-long
    LOG.info("Replay event log and verify re-calculated result against RTMR. Provide option for selected measurement verification.\n")

    parser = argparse.ArgumentParser(
        description="The utility to replay event logs logged at both boot time \
            and runtime, and verify the replayed results")
    parser.add_argument('--verify-register-index',
                        nargs="*", type=int,
                        help='list of register to verify', dest='index_info')
    parser.add_argument('-f', type=str, help='path to selected event log type',
                        dest='event_log_file')
    args = parser.parse_args()

    if args.index_info is not None and args.index_info != "":
        for i in args.index_info:
            if not 0 <= i < 4:
                LOG.error("Invalid RTMR registers specified. \
                          The value should be in range 0-3.")
                raise ValueError("Invalid RTMR value")
    else:
        args.index_info = [0, 1, 2, 3]

    # check if IMA over RTMR has been enabled
    IMA_RTMR_FLAG = True
    with open("/proc/cmdline", encoding="utf-8") as proc_f:
        cmdline = proc_f.read().splitlines()
        if "ima_hash=sha384" not in cmdline[0].split(" "):
            # pylint: disable-next=line-too-long
            LOG.info("IMA over RTMR not enabled. Verify basic boot measurements.")
            IMA_RTMR_FLAG = False

    actor = VerifyActor()
    actor.verify_rtmr(index_list=args.index_info,
                      event_log_file=args.event_log_file,
                      ima_flag=IMA_RTMR_FLAG)
