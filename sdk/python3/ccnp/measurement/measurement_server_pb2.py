# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: measurement-server.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x18measurement-server.proto\x12\x0bmeasurement\"\xa6\x01\n\x15GetMeasurementRequest\x12+\n\x10measurement_type\x18\x01 \x01(\x0e\x32\x11.measurement.TYPE\x12\x33\n\x14measurement_category\x18\x02 \x01(\x0e\x32\x15.measurement.CATEGORY\x12\x13\n\x0breport_data\x18\x03 \x01(\t\x12\x16\n\x0eregister_index\x18\x04 \x01(\x05\"*\n\x13GetMeasurementReply\x12\x13\n\x0bmeasurement\x18\x01 \x01(\t*\x1a\n\x04TYPE\x12\x08\n\x04PAAS\x10\x00\x12\x08\n\x04SAAS\x10\x01*1\n\x08\x43\x41TEGORY\x12\x0e\n\nTEE_REPORT\x10\x00\x12\x07\n\x03TPM\x10\x01\x12\x0c\n\x08TDX_RTMR\x10\x02\x32g\n\x0bMeasurement\x12X\n\x0eGetMeasurement\x12\".measurement.GetMeasurementRequest\x1a .measurement.GetMeasurementReply\"\x00\x42gZegithub.com/intel/confidential-cloud-native-primitives/service/measurement-server/proto/getMeasurementb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'measurement_server_pb2', _globals)
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Zegithub.com/intel/confidential-cloud-native-primitives/service/measurement-server/proto/getMeasurement'
  _globals['_TYPE']._serialized_start=254
  _globals['_TYPE']._serialized_end=280
  _globals['_CATEGORY']._serialized_start=282
  _globals['_CATEGORY']._serialized_end=331
  _globals['_GETMEASUREMENTREQUEST']._serialized_start=42
  _globals['_GETMEASUREMENTREQUEST']._serialized_end=208
  _globals['_GETMEASUREMENTREPLY']._serialized_start=210
  _globals['_GETMEASUREMENTREPLY']._serialized_end=252
  _globals['_MEASUREMENT']._serialized_start=333
  _globals['_MEASUREMENT']._serialized_end=436
# @@protoc_insertion_point(module_scope)
