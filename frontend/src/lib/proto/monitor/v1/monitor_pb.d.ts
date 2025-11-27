import * as jspb from 'google-protobuf'

import * as google_api_annotations_pb from '../../google/api/annotations_pb'; // proto import: "google/api/annotations.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb'; // proto import: "google/protobuf/empty.proto"


export class Device extends jspb.Message {
  getId(): string;
  setId(value: string): Device;

  getDeviceId(): string;
  setDeviceId(value: string): Device;

  getAlias(): string;
  setAlias(value: string): Device;

  getHost(): string;
  setHost(value: string): Device;

  getPort(): number;
  setPort(value: number): Device;

  getPortGateway(): number;
  setPortGateway(value: number): Device;

  getArchitecture(): string;
  setArchitecture(value: string): Device;

  getOs(): string;
  setOs(value: string): Device;

  getSupportedProtocolsList(): Array<Protocol>;
  setSupportedProtocolsList(value: Array<Protocol>): Device;
  clearSupportedProtocolsList(): Device;
  addSupportedProtocols(value: Protocol, index?: number): Device;

  getCreatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setCreatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Device;
  hasCreatedAt(): boolean;
  clearCreatedAt(): Device;

  getUpdatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: google_protobuf_timestamp_pb.Timestamp): Device;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): Device;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Device.AsObject;
  static toObject(includeInstance: boolean, msg: Device): Device.AsObject;
  static serializeBinaryToWriter(message: Device, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Device;
  static deserializeBinaryFromReader(message: Device, reader: jspb.BinaryReader): Device;
}

export namespace Device {
  export type AsObject = {
    id: string,
    deviceId: string,
    alias: string,
    host: string,
    port: number,
    portGateway: number,
    architecture: string,
    os: string,
    supportedProtocolsList: Array<Protocol>,
    createdAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    updatedAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class Diagnostics extends jspb.Message {
  getHardwareVersion(): string;
  setHardwareVersion(value: string): Diagnostics;

  getSoftwareVersion(): string;
  setSoftwareVersion(value: string): Diagnostics;

  getFirmwareVersion(): string;
  setFirmwareVersion(value: string): Diagnostics;

  getCpuUsage(): number;
  setCpuUsage(value: number): Diagnostics;

  getMemoryUsage(): number;
  setMemoryUsage(value: number): Diagnostics;

  getDeviceStatus(): DeviceStatus;
  setDeviceStatus(value: DeviceStatus): Diagnostics;

  getChecksum(): string;
  setChecksum(value: string): Diagnostics;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Diagnostics.AsObject;
  static toObject(includeInstance: boolean, msg: Diagnostics): Diagnostics.AsObject;
  static serializeBinaryToWriter(message: Diagnostics, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Diagnostics;
  static deserializeBinaryFromReader(message: Diagnostics, reader: jspb.BinaryReader): Diagnostics;
}

export namespace Diagnostics {
  export type AsObject = {
    hardwareVersion: string,
    softwareVersion: string,
    firmwareVersion: string,
    cpuUsage: number,
    memoryUsage: number,
    deviceStatus: DeviceStatus,
    checksum: string,
  }
}

export class RegisterDeviceRequest extends jspb.Message {
  getDeviceId(): string;
  setDeviceId(value: string): RegisterDeviceRequest;

  getAlias(): string;
  setAlias(value: string): RegisterDeviceRequest;

  getHost(): string;
  setHost(value: string): RegisterDeviceRequest;

  getPort(): number;
  setPort(value: number): RegisterDeviceRequest;

  getPortGateway(): number;
  setPortGateway(value: number): RegisterDeviceRequest;

  getProtocol(): Protocol;
  setProtocol(value: Protocol): RegisterDeviceRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterDeviceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterDeviceRequest): RegisterDeviceRequest.AsObject;
  static serializeBinaryToWriter(message: RegisterDeviceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterDeviceRequest;
  static deserializeBinaryFromReader(message: RegisterDeviceRequest, reader: jspb.BinaryReader): RegisterDeviceRequest;
}

export namespace RegisterDeviceRequest {
  export type AsObject = {
    deviceId: string,
    alias: string,
    host: string,
    port: number,
    portGateway: number,
    protocol: Protocol,
  }
}

export class RegisterDeviceResponse extends jspb.Message {
  getDevice(): Device | undefined;
  setDevice(value?: Device): RegisterDeviceResponse;
  hasDevice(): boolean;
  clearDevice(): RegisterDeviceResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterDeviceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterDeviceResponse): RegisterDeviceResponse.AsObject;
  static serializeBinaryToWriter(message: RegisterDeviceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterDeviceResponse;
  static deserializeBinaryFromReader(message: RegisterDeviceResponse, reader: jspb.BinaryReader): RegisterDeviceResponse;
}

export namespace RegisterDeviceResponse {
  export type AsObject = {
    device?: Device.AsObject,
  }
}

export class ListDevicesResponse extends jspb.Message {
  getDevicesList(): Array<Device>;
  setDevicesList(value: Array<Device>): ListDevicesResponse;
  clearDevicesList(): ListDevicesResponse;
  addDevices(value?: Device, index?: number): Device;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListDevicesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListDevicesResponse): ListDevicesResponse.AsObject;
  static serializeBinaryToWriter(message: ListDevicesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListDevicesResponse;
  static deserializeBinaryFromReader(message: ListDevicesResponse, reader: jspb.BinaryReader): ListDevicesResponse;
}

export namespace ListDevicesResponse {
  export type AsObject = {
    devicesList: Array<Device.AsObject>,
  }
}

export class UpdateDeviceRequest extends jspb.Message {
  getDeviceId(): string;
  setDeviceId(value: string): UpdateDeviceRequest;

  getDeviceStatus(): DeviceStatus;
  setDeviceStatus(value: DeviceStatus): UpdateDeviceRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateDeviceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateDeviceRequest): UpdateDeviceRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateDeviceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateDeviceRequest;
  static deserializeBinaryFromReader(message: UpdateDeviceRequest, reader: jspb.BinaryReader): UpdateDeviceRequest;
}

export namespace UpdateDeviceRequest {
  export type AsObject = {
    deviceId: string,
    deviceStatus: DeviceStatus,
  }
}

export class DiagnosticsRequest extends jspb.Message {
  getDeviceId(): string;
  setDeviceId(value: string): DiagnosticsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DiagnosticsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DiagnosticsRequest): DiagnosticsRequest.AsObject;
  static serializeBinaryToWriter(message: DiagnosticsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DiagnosticsRequest;
  static deserializeBinaryFromReader(message: DiagnosticsRequest, reader: jspb.BinaryReader): DiagnosticsRequest;
}

export namespace DiagnosticsRequest {
  export type AsObject = {
    deviceId: string,
  }
}

export class DiagnosticsResponse extends jspb.Message {
  getDevice(): Device | undefined;
  setDevice(value?: Device): DiagnosticsResponse;
  hasDevice(): boolean;
  clearDevice(): DiagnosticsResponse;

  getDiagnostics(): Diagnostics | undefined;
  setDiagnostics(value?: Diagnostics): DiagnosticsResponse;
  hasDiagnostics(): boolean;
  clearDiagnostics(): DiagnosticsResponse;

  getUpdatedAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setUpdatedAt(value?: google_protobuf_timestamp_pb.Timestamp): DiagnosticsResponse;
  hasUpdatedAt(): boolean;
  clearUpdatedAt(): DiagnosticsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DiagnosticsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DiagnosticsResponse): DiagnosticsResponse.AsObject;
  static serializeBinaryToWriter(message: DiagnosticsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DiagnosticsResponse;
  static deserializeBinaryFromReader(message: DiagnosticsResponse, reader: jspb.BinaryReader): DiagnosticsResponse;
}

export namespace DiagnosticsResponse {
  export type AsObject = {
    device?: Device.AsObject,
    diagnostics?: Diagnostics.AsObject,
    updatedAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export enum Protocol { 
  PROTOCOL_UNSPECIFIED = 0,
  PROTOCOL_HTTP = 1,
  PROTOCOL_HTTP_STREAM = 2,
  PROTOCOL_GRPC = 3,
  PROTOCOL_GRPC_STREAM = 4,
}
export enum DeviceStatus { 
  DEVICE_STATUS_UNSPECIFIED = 0,
  DEVICE_STATUS_HEALTHY = 1,
  DEVICE_STATUS_DEGRADED = 2,
  DEVICE_STATUS_ERROR = 3,
  DEVICE_STATUS_MAINTENANCE = 4,
  DEVICE_STATUS_BOOTING = 5,
  DEVICE_STATUS_OFFLINE = 6,
}
