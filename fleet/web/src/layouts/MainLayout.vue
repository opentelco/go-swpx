<template>
  <q-layout view="hHh lpR fFf">
    <q-header elevated>
      <q-toolbar class="">
        <q-btn flat to="/">
          <img src="~assets/ot-logo-white.png" class="toolbar-logo" />
        </q-btn>
        <q-space />

        <q-btn-dropdown flat label="Fleet">

          <div class="row no-wrap q-pa-md">
            <q-item clickable tag="a" to="/devices" class="mainmenu-item">
              <q-item-section avatar>
                <q-icon name="fa-solid fa-network-wired" />
              </q-item-section>

              <q-item-section>
                <q-item-label>Devices</q-item-label>
                <q-item-label caption>List and manage your fleet of devices</q-item-label>
              </q-item-section>
            </q-item>

            <q-item clickable class="mainmenu-item" disabled>
              <q-item-section avatar>
                <q-icon name="fa-solid fa-ethernet" />
              </q-item-section>

              <q-item-section class="col-12">
                <q-item-label>Ports</q-item-label>
                <q-item-label caption>List all ports connected to devices</q-item-label>
                <q-item-label><q-badge color="secondary" align="middle">COMING SOON</q-badge></q-item-label>
              </q-item-section>
            </q-item>


          </div>
          <div class="row no-wrap q-pa-md">
            <q-item clickable tag="a" to="/devices/events" class="mainmenu-item">
              <q-item-section avatar>
                <q-icon name="fa-solid fa-envelope-open-text" />
              </q-item-section>

              <q-item-section>
                <q-item-label>Events</q-item-label>
                <q-item-label caption>Events occured on any device in the fleet</q-item-label>
              </q-item-section>
            </q-item>
            <q-item clickable tag="a" to="/devices/changes" class="mainmenu-item">
              <q-item-section avatar>
                <q-icon name="fa-solid fa-code-pull-request" />
              </q-item-section>

              <q-item-section>
                <q-item-label>Changes</q-item-label>
                <q-item-label caption>Changes made on any device in the fleet</q-item-label>
              </q-item-section>
            </q-item>
          </div>
        </q-btn-dropdown>

        <q-btn-dropdown flat label="Config">
          <div class="row no-wrap q-pa-md">
            <q-item clickable tag="a" to="/swp" class="mainmenu-item">
              <q-item-section avatar>
                <q-icon name="data_object" />
              </q-item-section>

              <q-item-section>
                <q-item-label>Stanza</q-item-label>
                <q-item-label caption>Stanzas (configuration snippets) for your devices</q-item-label>
              </q-item-section>
            </q-item>

            <q-item clickable tag="a" to="/devices/configuration" class="mainmenu-item">
              <q-item-section avatar>
                <q-icon name="fa-solid fa-file-code" />
              </q-item-section>

              <q-item-section>
                <q-item-label>Configuration</q-item-label>
                <q-item-label caption>List running configs for devices</q-item-label>
              </q-item-section>
            </q-item>


          </div>
          <div class="row no-wrap q-pa-md">
            <q-item clickable class="mainmenu-item" disabled>
              <q-item-section avatar>
                <q-icon name="fa-solid fa-book" />
              </q-item-section>

              <q-item-section>
                <q-item-label>Templates</q-item-label>
                <q-item-label caption>Configuration templates for your devices</q-item-label>
                <q-item-label><q-badge color="secondary" align="middle">COMING SOON</q-badge></q-item-label>
              </q-item-section>
            </q-item>


            <q-item clickable class="mainmenu-item" disabled>
              <q-item-section avatar>
                <q-icon name="bolt" />
              </q-item-section>

              <q-item-section>
                <q-item-label>Switch Poller</q-item-label>
                <q-item-label caption>use the poller to get data from your devices</q-item-label>
                <q-item-label><q-badge color="secondary" align="middle">COMING SOON</q-badge></q-item-label>
              </q-item-section>
            </q-item>

          </div>
        </q-btn-dropdown>

        <q-btn-dropdown flat label="Admin">
          <div class="row no-wrap q-pa-md">
            <q-item clickable tag="a" to="/admin/notifications" class="mainmenu-item">
              <q-item-section avatar>
                <q-icon name="fa-solid fa-bell" />
              </q-item-section>

              <q-item-section>
                <q-item-label>Notifications</q-item-label>
                <q-item-label caption>List and manage notifications</q-item-label>
              </q-item-section>
            </q-item>

            <q-item clickable class="mainmenu-item" disabled>
              <q-item-section avatar>
                <q-icon name="fa-solid fa-copyright" />
              </q-item-section>

              <q-item-section class="col-12">
                <q-item-label>Vendors</q-item-label>
                <q-item-label caption>Manage available vendors</q-item-label>
                <q-item-label><q-badge color="secondary" align="middle">COMING SOON</q-badge></q-item-label>
              </q-item-section>
            </q-item>


          </div>
          <div class="row no-wrap q-pa-md justify-end">

            <q-item clickable tag="a" to="/admin" class="mainmenu-item">
              <q-item-section avatar>
                <q-icon name="fa-solid fa-ellipsis" />
              </q-item-section>

              <q-item-section>
                <q-item-label>Administration</q-item-label>
                <q-item-label caption>Go to the settings page for all knobs and settings </q-item-label>
              </q-item-section>
            </q-item>
          </div>
        </q-btn-dropdown>

        <q-btn flat dense round icon="dark_mode" @click="toggleDarMode" />
        <NotificationsMainButton />

        <q-item-section avatar>
          <q-avatar>
            <img src="https://cdn.quasar.dev/img/avatar3.jpg" />
          </q-avatar>
        </q-item-section>
      </q-toolbar>
    </q-header>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';
import { useQuasar } from 'quasar';
import NotificationsMainButton from 'src/components/NotificationsMainButton.vue';

export default defineComponent({
  name: 'MainLayout',

  components: {
    NotificationsMainButton,
  },

  setup() {
    const leftDrawerOpen = ref(false);

    const $q = useQuasar()

    // calling here; equivalent to when component is created


    return {
      leftDrawerOpen,
      toggleLeftDrawer() {
        leftDrawerOpen.value = !leftDrawerOpen.value;
      },
      toggleDarMode() {
        $q.dark.toggle()
      }
    };
  },
});
</script>
