<script setup lang="ts">

import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'

import { ListNotificationsResponse, Notification, PageInfo } from 'src/gql/graphql';

interface notifications {
  notifications: ListNotificationsResponse

}

const { result, loading, error } = useQuery<notifications>(gql`
  query Notifications ($limit: Int!){
    notifications(params: { limit: $limit } ) {
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
  limit: 0,
})


function timeDifference(current, previous) {
  previous = new Date(previous)

  var msPerMinute = 60 * 1000;
  var msPerHour = msPerMinute * 60;
  var msPerDay = msPerHour * 24;
  var msPerMonth = msPerDay * 30;
  var msPerYear = msPerDay * 365;

  var elapsed = current - previous;

  if (elapsed < msPerMinute) {
    return Math.round(elapsed / 1000) + ' seconds ago';
  }

  else if (elapsed < msPerHour) {
    return Math.round(elapsed / msPerMinute) + ' minutes ago';
  }

  else if (elapsed < msPerDay) {
    return Math.round(elapsed / msPerHour) + ' hours ago';
  }

  else if (elapsed < msPerMonth) {
    return '~' + Math.round(elapsed / msPerDay) + ' days ago';
  }

  else if (elapsed < msPerYear) {
    return '~' + Math.round(elapsed / msPerMonth) + ' months ago';
  }

  else {
    return '~' + Math.round(elapsed / msPerYear) + ' years ago';
  }
}
</script>

<template>
  <q-btn flat dense round icon="notifications" class="q-mr-md">
    <q-badge color="red" floating v-if="!error && result?.notifications.notifications" round>
      {{ result.notifications.pageInfo.count }}
    </q-badge>
    <q-menu>
      <q-list bordered class="rounded-borders" style="max-width: 400px" v-if="result && result.notifications.notifications">
        <q-item-label header>Notifications</q-item-label>

        <q-item clickable v-ripple v-for="n in result.notifications.notifications" :key="n.id" >
          <q-item-section avatar>
            <q-avatar>
              <img src="https://cdn.quasar.dev/img/avatar2.jpg">
            </q-avatar>
          </q-item-section>

          <q-item-section>
            <q-item-label lines="1">{{ n.title }}</q-item-label>
            <q-item-label caption lines="3">
              {{ n.message }}
            </q-item-label>
          </q-item-section>

          <q-item-section side top>
            {{ timeDifference(new Date(), n.timestamp) }}
          </q-item-section>
        </q-item>


      </q-list>
    </q-menu>
  </q-btn>
</template>
