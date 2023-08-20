/* eslint-disable */
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Duration: { input: any; output: any; }
  Timestamp: { input: any; output: any; }
};

export type ChangeConnection = {
  __typename?: 'ChangeConnection';
  changes?: Maybe<Array<DeviceChange>>;
  pageInfo: PageInfo;
};

export type Configuration = {
  __typename?: 'Configuration';
  changes?: Maybe<Scalars['String']['output']>;
  checksum?: Maybe<Scalars['String']['output']>;
  configuration?: Maybe<Scalars['String']['output']>;
  createdAt: Scalars['Timestamp']['output'];
  device?: Maybe<Device>;
  id: Scalars['ID']['output'];
};

export type ConfigurationConnection = {
  __typename?: 'ConfigurationConnection';
  configurations?: Maybe<Array<Configuration>>;
  pageInfo: PageInfo;
};

export type Device = {
  __typename?: 'Device';
  changes: ChangeConnection;
  configurations: ConfigurationConnection;
  createdAt: Scalars['Timestamp']['output'];
  domain?: Maybe<Scalars['String']['output']>;
  events: EventConnection;
  hostname?: Maybe<Scalars['String']['output']>;
  id: Scalars['ID']['output'];
  lastReboot?: Maybe<Scalars['Timestamp']['output']>;
  lastSeen?: Maybe<Scalars['Timestamp']['output']>;
  managementIp?: Maybe<Scalars['String']['output']>;
  model?: Maybe<Scalars['String']['output']>;
  networkRegion?: Maybe<Scalars['String']['output']>;
  pollerProviderPlugin?: Maybe<Scalars['String']['output']>;
  pollerResourcePlugin?: Maybe<Scalars['String']['output']>;
  schedules?: Maybe<Array<DeviceSchedule>>;
  serialNumber?: Maybe<Scalars['String']['output']>;
  stanzas: StanzaConnection;
  state: DeviceState;
  status: DeviceStatus;
  updatedAt: Scalars['Timestamp']['output'];
  version?: Maybe<Scalars['String']['output']>;
};


export type DeviceChangesArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
};


export type DeviceConfigurationsArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
};


export type DeviceEventsArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
};


export type DeviceStanzasArgs = {
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
};

export type DeviceChange = {
  __typename?: 'DeviceChange';
  createdAt: Scalars['Timestamp']['output'];
  field: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  newValue: Scalars['String']['output'];
  oldValue: Scalars['String']['output'];
};

export type DeviceEvent = {
  __typename?: 'DeviceEvent';
  action: DeviceEventAction;
  createdAt: Scalars['Timestamp']['output'];
  id: Scalars['ID']['output'];
  message: Scalars['String']['output'];
  outcome: DeviceEventOutcome;
  type: DeviceEventType;
};

export enum DeviceEventAction {
  CollectConfig = 'COLLECT_CONFIG',
  CollectDevice = 'COLLECT_DEVICE',
  Create = 'CREATE',
  NotSet = 'NOT_SET',
  Update = 'UPDATE'
}

export enum DeviceEventOutcome {
  Failure = 'FAILURE',
  NotSet = 'NOT_SET',
  Success = 'SUCCESS'
}

export enum DeviceEventType {
  Configuration = 'CONFIGURATION',
  Device = 'DEVICE',
  NotSet = 'NOT_SET'
}

export type DeviceSchedule = {
  __typename?: 'DeviceSchedule';
  active: Scalars['Boolean']['output'];
  failedCount: Scalars['Int']['output'];
  interval: Scalars['Duration']['output'];
  lastRun?: Maybe<Scalars['Timestamp']['output']>;
  type: ScheduleType;
};

export enum DeviceState {
  Active = 'ACTIVE',
  Inactive = 'INACTIVE',
  New = 'NEW',
  NotSet = 'NOT_SET',
  Rouge = 'ROUGE'
}

export enum DeviceStatus {
  New = 'NEW',
  NotSet = 'NOT_SET',
  Reachable = 'REACHABLE',
  Unreachable = 'UNREACHABLE'
}

export type EventConnection = {
  __typename?: 'EventConnection';
  events?: Maybe<Array<DeviceEvent>>;
  pageInfo: PageInfo;
};

export type ListDeviceResponse = {
  __typename?: 'ListDeviceResponse';
  devices?: Maybe<Array<Device>>;
  pageInfo: PageInfo;
};

export type ListDevicesParams = {
  hostname?: InputMaybe<Scalars['String']['input']>;
  limit?: InputMaybe<Scalars['Int']['input']>;
  managementIp?: InputMaybe<Scalars['String']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
  search?: InputMaybe<Scalars['String']['input']>;
};

export enum ListNotificationFilter {
  IncludeRead = 'INCLUDE_READ',
  ResourceTypeConfig = 'RESOURCE_TYPE_CONFIG',
  ResourceTypeDevice = 'RESOURCE_TYPE_DEVICE'
}

export type ListNotificationsParams = {
  filter?: InputMaybe<Array<ListNotificationFilter>>;
  limit?: InputMaybe<Scalars['Int']['input']>;
  offset?: InputMaybe<Scalars['Int']['input']>;
  resource_ids?: InputMaybe<Array<Scalars['ID']['input']>>;
};

export type ListNotificationsResponse = {
  __typename?: 'ListNotificationsResponse';
  notifications: Array<Notification>;
  pageInfo: PageInfo;
};

export type MarkNotificationsAsReadParams = {
  ids: Array<Scalars['ID']['input']>;
};

export type Mutation = {
  __typename?: 'Mutation';
  markNotificationsAsRead?: Maybe<Array<Notification>>;
};


export type MutationMarkNotificationsAsReadArgs = {
  input: MarkNotificationsAsReadParams;
};

export type Notification = {
  __typename?: 'Notification';
  id: Scalars['ID']['output'];
  message: Scalars['String']['output'];
  read: Scalars['Boolean']['output'];
  resource_id: Scalars['ID']['output'];
  resource_type: NotificationResourceType;
  timestamp: Scalars['Timestamp']['output'];
  title: Scalars['String']['output'];
};

export enum NotificationResourceType {
  Config = 'CONFIG',
  Device = 'DEVICE',
  Unspecified = 'UNSPECIFIED'
}

export type PageInfo = {
  __typename?: 'PageInfo';
  count?: Maybe<Scalars['Int']['output']>;
  limit?: Maybe<Scalars['Int']['output']>;
  offset?: Maybe<Scalars['Int']['output']>;
  total?: Maybe<Scalars['Int']['output']>;
};

export type Query = {
  __typename?: 'Query';
  device: Device;
  devices: ListDeviceResponse;
  notifications: ListNotificationsResponse;
};


export type QueryDeviceArgs = {
  id: Scalars['ID']['input'];
};


export type QueryDevicesArgs = {
  params?: InputMaybe<ListDevicesParams>;
};


export type QueryNotificationsArgs = {
  params?: InputMaybe<ListNotificationsParams>;
};

export enum ScheduleType {
  Config = 'CONFIG',
  Device = 'DEVICE',
  NotSet = 'NOT_SET'
}

export type Stanza = {
  __typename?: 'Stanza';
  appliedAt: Scalars['Timestamp']['output'];
  content?: Maybe<Scalars['String']['output']>;
  createdAt: Scalars['Timestamp']['output'];
  description?: Maybe<Scalars['String']['output']>;
  device?: Maybe<Device>;
  device_type?: Maybe<Scalars['String']['output']>;
  id: Scalars['ID']['output'];
  name?: Maybe<Scalars['String']['output']>;
  revert_content?: Maybe<Scalars['String']['output']>;
  revert_template?: Maybe<Scalars['String']['output']>;
  template?: Maybe<Scalars['String']['output']>;
  updatedAt: Scalars['Timestamp']['output'];
};

export type StanzaConnection = {
  __typename?: 'StanzaConnection';
  pageInfo: PageInfo;
  stanzas?: Maybe<Array<Stanza>>;
};
