# Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
# SPDX-License-Identifier: Apache-2.0

"""
This package provides the definitions and helper class for Quote of confidetial computing,
which will be used for remote attestation.

Reference:
1. Part 2: Structures, Trusted Platform Module Library
https://trustedcomputinggroup.org/wp-content/uploads/TCG_TPM2_r1p59_Part2_Structures_pub.pdf
2. Architecture Specification: Intel® Trust Domain Extensions (Intel® TDX) Module
https://cdrdv2.intel.com/v1/dl/getContent/733568
"""

import base64
import logging
import os
import struct
from typing import Optional
import grpc
# pylint: disable=E1101
from ccnp.quote import quote_server_pb2
from ccnp.quote import quote_server_pb2_grpc

LOG = logging.getLogger(__name__)

# Default gRPC timeout
TIMEOUT = 5

class QuoteClient:
    """Quote client class

    This class is a client to connect to Quote Server and do gRPC call getting the
    server.

    Attributes:
        _server (str): The gRPC server to connect.
        _channel (Channel): The gRPC channel, thread-safe.
        _nonce (str): The nonce parameter to get quote.
        _user_data (str): The user data parameter to get quote.
        _stub (GetQuoteStub): The get quote stub for gRPC.
        _request (GetQuoteRequest): The get quote request for gRPC.
    """
    def __init__(self, server: str="unix:/run/ccnp/uds/quote-server.sock"):
        """Initialize a quote client object

        This constructor initializes quote client object with Unix Domain Socket (UDS)
        path. And prepare default atrributes.

        Args:
            server (str): gRPC server UDS path, default is /run/ccnp/uds/quote-server.sock
        
        Raises:
            ValueError: If server UDS path is not valid.
        """
        if len(server) == 0 or server[:5] != "unix:":
            raise ValueError("Invalid server path, only unix domain socket supported.")
        self._server = server
        self._channel = None
        self._stub = None
        self._request = None

    def request(self, nonce: str, user_data: str) -> Optional[quote_server_pb2.GetQuoteResponse]:
        """Do reuqest to Quote Server
        Detect the Quote Server and gRPC connect to the server. Make the getting quote stub
        and request for communication.

        Args:
            nonce (str): The nonce parameters for getting quote.
            user_data (str): The user data parameters for getting quote.

        Raises:
            RuntimeError: If Quote Server does not start.
            ConnectionRefusedError: If connect to Quote Server failed.
        """
        if self._channel is None:
            if not os.path.exists(self._server.replace('unix:', '')):
                raise RuntimeError("Quote server does not start.")
            self._channel = grpc.insecure_channel(self._server,
                                                  options=[('grpc.default_authority', 'localhost')])
            try:
                grpc.channel_ready_future(self._channel).result(timeout=TIMEOUT)
            except grpc.FutureTimeoutError as err:
                raise ConnectionRefusedError('Connection to quote server failed') from err
            self._stub = quote_server_pb2_grpc.GetQuoteStub(self._channel)

        self._request = quote_server_pb2.GetQuoteRequest(nonce=nonce, user_data=user_data)
        resp = self._stub.GetQuote(self._request)
        return resp

class Quote():
    """An abstract base class for Quote

    This class a abstract class with a common static method `get_quote` for external
    SDK interface, the subclasses need to implement `parse` method to parse Quote
    information.

    Attributes:
        _quote (bytes): The bytes of a quote.
        _type (str): The type of a quote.
    """

    TYPE_TDX = 'TDX'
    TYPE_TPM = 'TPM'

    def __init__(self, quote: str = None, quote_type: str = None):
        """Initialize Quote object.

        The constructor to initialize quote object with quote bytes and quote type.

        Args:
            quote (bytes): The bytes of a quote.
            quote_type (str): The type of a quote.
        """
        self._quote = quote
        self._type = quote_type

    @property
    def quote_type(self) -> int:
        """str: The type of the quote."""
        return self._type

    @property
    def quote(self) -> bytes:
        """bytes: the bytes of the quote"""
        return self._quote

    @staticmethod
    def get_quote(nonce: str = None, user_data: str = None):
        """Get quote interface

        The get quote interface to expose to SDK.

        Args:
            nonce (str): Base64 encoded nonce to prevent replay attack.
            user_data (str): Base64 encoded user data to be measured in a quote.

        Returns:
            Quote: The quote object for specific quote type.
            None: Filed to get a quote.
        """
        quote = QuoteClient()
        resp = quote.request(nonce, user_data)
        if resp is not None and resp.quote_type is not None \
            and resp.quote is not None and len(resp.quote) > 0:
            LOG.info("Get quote successfully.")
            quote_data = base64.b64decode(resp.quote)
            quote_type = resp.quote_type
            if quote_type == "TDX":
                td_quote = QuoteTDX(quote_data, quote_type)
                td_quote.parse()
                return td_quote
        LOG.error("Failed to get quote.")
        return None

class QuoteTDX(Quote):
    """TDX quote class

    This class is a subclass of Quote to parse TDX sepecific quote.
    Refer: https://cdrdv2.intel.com/v1/dl/getContent/733568

    Attributes:
        _version (int): TD quote version
        _tdreport (bytes): The bytes of TD report.
        _tee_type (int): Type of TEE for which the Quote has been generated.
        _tee_tcb_svn (bytes): Array of TEE TCB SVNs.
        _mrseam (bytes): Measurement of the SEAM module (SHA384 hash). 
        _mrsignerseam (bytes): Measurement of a 3rd party SEAM module’s signer (SHA384 hash).
        _seamattributes (bytes): SEAM’s ATTRIBUTES.
        _tdattributes (bytes): TD’s ATTRIBUTES.
        _xfam (bytes): TD’s XFAM.
        _mrtd (bytes): Measurement of the initial contents of the TD (SHA384 hash).
        _mrconfigid (bytes): Software defined ID for non-owner-defined configuration of the TD
        _mrowner (bytes): Software defined ID for the guest TD’s owner.
        _mrownerconfig (bytes): Software defined ID for owner-defined configuration of the TD
        _rtmr (bytes): Array of 4 runtime extendable measurement registers (SHA384 hash).
        _reportdata (bytes): Additional Report Data.
        _signature (bytes): ECDSA signature, r component followed by s component, 2 x 32 bytes.
        _attestation_key (bytes): Public part of ECDSA Attestation Key generated by Quoting Enclave.
        _cert_data (bytes): Data required to certify Attestation Key used to sign the Quote.
    """
    def __init__(self, quote: bytes, quote_type: str):
        """Initialize TD quote object

        The constructor of TD quote object, initialize attributes.

        Args:
            quote (bytes): The bytes of a quote.
            quote_type (str): The type of a quote.
        """
        super().__init__(quote, quote_type)
        self._version = 0
        self._tdreport = None
        self._tee_type = 0
        self._tee_tcb_svn = None
        self._mrseam = None
        self._mrsignerseam = None
        self._seamattributes = None
        self._tdattributes = None
        self._xfam = None
        self._mrtd = None
        self._mrconfigid = None
        self._mrowner = None
        self._mrownerconfig = None
        self._rtmrs = []
        self._reportdata = None
        self._signature = None
        self._attestation_key = None
        self._cert_data = None

    @property
    def version(self) -> int:
        """int: the version of the quote"""
        return self._version

    @property
    def tdreport(self) -> bytes:
        """bytes: the bytes of the TD report"""
        return self._tdreport

    @property
    def tee_type(self) -> int:
        """int: the TEE type of the quote"""
        return self._tee_type

    @property
    def mrseam(self) -> bytes:
        """bytes: the MRSEAM in the quote"""
        return self._mrseam

    @property
    def mrsignerseam(self) -> bytes:
        """bytes: the bytes of MRSIGNERSEAM in the quote"""
        return self._mrsignerseam

    @property
    def seam_attributes(self) -> bytes:
        """bytes: the bytes of SEAM ATTRIBUTES in the quote"""
        return self._seamattributes

    @property
    def td_attributes(self) -> bytes:
        """bytes: the bytes of TD ATTRIBUTES in the quote"""
        return self._tdattributes

    @property
    def xfam(self) -> bytes:
        """bytes: the bytes of XFAM in the quote"""
        return self._xfam

    @property
    def mrtd(self) -> bytes:
        """bytes: the bytes of MRTD in the quote"""
        return self._mrtd

    @property
    def mrconfigid(self) -> bytes:
        """bytes: the bytes of MRCONFIGID in the quote"""
        return self._mrconfigid

    @property
    def mrowner(self) -> bytes:
        """bytes: the bytes of MROWNER in the quote"""
        return self._mrowner

    @property
    def mrownerconfig(self) -> bytes:
        """bytes: the bytes of MROWNERCONFIG in the quote"""
        return self._mrownerconfig

    @property
    def rtmrs(self) -> bytes:
        """bytes: the bytes of RTMRs in the quote"""
        rtmrs=[]
        for i in range(4):
            rtmrs.append(self._rtmrs[i*48:(i+1)*48])
        return rtmrs

    @property
    def report_data(self) -> bytes:
        """bytes: the bytes of REPORTDATA in the quote"""
        return self._reportdata

    @property
    def signature(self) -> bytes:
        """bytes: the bytes of signature of the quote"""
        return self._signature

    @property
    def attestation_key(self) -> bytes:
        """bytes: the bytes of attestation key in the quote"""
        return self._attestation_key

    @property
    def cert_data(self) -> bytes:
        """bytes: the bytes of certification data in the quote"""
        return self._cert_data

    def parse(self):
        """Parse TD quote

        This method is to parse the TD quote and TD report data.
        Refer: https://cdrdv2.intel.com/v1/dl/getContent/733568

        Raises:
            struct.error: Unpack quote data failed.
        """
        quote_len = len(self._quote)
        # Header, Body, Auth Data
        header, self._tdreport, auth_size, auth_data = \
            struct.unpack(f"<48s584sI{quote_len-48-584-4}s", self._quote)
        auth_data = auth_data[:auth_size]
        # Header
        self._version, _, self._tee_type, _ = struct.unpack(f"<2HI{len(header)-8}s", header)
        # Body
        self._tee_tcb_svn, self._mrseam, self._mrsignerseam, self._seamattributes, \
        self._tdattributes, self._xfam, self._mrtd, self._mrconfigid, self._mrowner, \
        self._mrownerconfig, self._rtmrs, self._reportdata = \
            struct.unpack("16s48s48s8s8s8s48s48s48s48s192s64s", self._tdreport)
        # Auth Data
        self._signature, self._attestation_key, cert_data = \
            struct.unpack(f"64s64s{auth_size-128}s", auth_data)
        # Certification Data
        _, _, cert_data = struct.unpack(f"<HI{len(cert_data)-6}s", cert_data)

class QuoteTPM(Quote):
    """TODO: implement TPM Quote class"""
    def __init__(self, quote: bytes, quote_type: str):
        """Initialize TPM quote object

        The constructor of TD quote object, initialize attributes.

        Args:
            quote (bytes): The bytes of a quote.
            quote_type (str): The type of a quote.
        """
        super().__init__(quote, quote_type)
        # TODO
