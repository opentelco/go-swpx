<script setup lang="ts">
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'

import { ref } from 'vue'
import { Notification, ListNotificationsResponse, ListNotificationFilter, NotificationResourceType,  } from 'src/gql/graphql';

interface notifications {
  notifications: ListNotificationsResponse

}
let defaultPageSize = ref(10)
let includeRead = ref(true)
let filters: ListNotificationFilter[] = [ListNotificationFilter.IncludeRead]

const { result, loading, error, refetch } = useQuery<notifications>(gql`
  query NotificationsList ($limit: Int!, $filter: [ListNotificationFilter!]){
    notifications(params: { limit: $limit, filter: $filter } ) {
        notifications {
            id
            title
            resource_id
            resource_type
            timestamp
            message
            read
        }
        pageInfo {
            limit
            offset
            total
            count
        }
      }
    }
`, {
  limit: defaultPageSize,
  filter: filters,
}, {})

const toggleReadFilter = () => {
  if (!includeRead.value) {
    filters = filters.filter((f) => f !== ListNotificationFilter.IncludeRead)
  } else {
    filters.push(ListNotificationFilter.IncludeRead)
  }
  refetch({
    limit: defaultPageSize.value,
    filter: filters,})

}
// typeToIcon converts a resource type to a material icon name
const typeToIcon = (type: NotificationResourceType) => {
  switch (type) {
    case NotificationResourceType.Device:
      return 'fa-solid fa-network-wired'
    case NotificationResourceType.Config:
      return 'fa-solid fa-file-code'
  }
}



const columns = [
  { name: 'read', label: '', field: '', align: 'center', },
  { name: 'title', label: 'Title', field: 'title',  align: 'left', },
  { name: 'message', label: 'Message', field: 'message', align: 'left', },
  { name: 'timestamp', label: 'Timestamp', field: 'timestamp', align: 'left', },
  { name: 'tools', label: '', field: '', align: 'center', },

]

const goTo = (notification: Notification) => {
  switch (notification.resource_type) {
    case NotificationResourceType.Device:
      return `/devices/${notification.resource_id}`
    case NotificationResourceType.Config:
      return `/config/${notification.resource_id}`
  }
}

</script>

<template>
  <div class="q-pa-md">
    <q-table :grid="$q.screen.xs" flat bordered title="Notifications" :rows="result?.notifications.notifications"
      :columns="columns" row-key="name" wrap-cells :loading="loading">

      <template v-slot:top>
        <span class="text-h6">Notifications</span>
        <q-space />
        <q-toggle v-model="includeRead" color="pink" icon="mark_chat_read" unchecked-icon="mark_chat_unread" label="Include read notifications" @click="toggleReadFilter"/>


      </template>

      <template v-slot:body="props">
        <q-tr :props="props">

          <q-td key="read" :props="props" class="vertical-top">
            <q-icon name="circle" size="1em" color="red" class="q-mx-md" v-if="!props.row.read" />
            <q-icon :name="typeToIcon(props.row.resource_type)" />
          </q-td>

          <q-td key="resource_type" :props="props" class="vertical-top">

          </q-td>

          <q-td key="title" :props="props" class="text-body1 vertical-top">
            {{ props.row.title }}
          </q-td>
          <q-td key="message" :props="props" class="text-body1 vertical-top" style="max-width:400px">
            {{ props.row.message }}
          </q-td>
          <q-td key="timestamp" :props="props" class="text-caption vertical-top">
            {{ props.row.timestamp }}
          </q-td>
          <q-td  key="tools" :props="props" class="vertical-top">
            <q-btn flat dense round icon="mark_chat_read" color="primary" class="" v-if="!props.row.read">
              <q-tooltip class="bg-primary">Mark notifications as read</q-tooltip>
            </q-btn>

            <q-btn flat dense round icon="link" color="secondary" class="" :to="goTo(props.row)">
              <q-tooltip class="bg-primary">Go to resource</q-tooltip>
            </q-btn>
          </q-td>


        </q-tr>
      </template>
    </q-table>

  </div>
</template>
