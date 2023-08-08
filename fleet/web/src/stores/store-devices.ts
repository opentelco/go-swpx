import { defineStore } from 'pinia';

import { ListDeviceResponse } from '../gql/graphql'

enum LoadingState {
  GetDevices,
}

export const useDevicesStore = defineStore('devices', {
  state: () => ({
    devices: {} as ListDeviceResponse,
    loading: [] as Array<LoadingState>,

  }),
  getters: {
    loading: (state) => state.loading.length > 0,
    getDevices: (state) => state.devices,
  },
  actions: {
    async fetchDevices() {
      this.loading.push(LoadingState.GetDevices);

    },
  },
});
