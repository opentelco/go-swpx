<template>
  <div class="q-pa-md">
    <div class="row">
      <div class="col timeline-options">
        <button @click="setGrouping('day')">day</button>
        <button @click="setGrouping('3day')">3day</button>
        <button @click="setGrouping('month')">month</button>
        <button @click="setGrouping('quarter')">quarter</button>
        <button @click="setGrouping('year')">1 Year</button>
      </div>
    </div>
    <div class="row">
      <div class="row timeline">
        <div v-for="(item, index) in groupedItems" :key="index" class="timeline-item col">
          <q-card class="shadow-3 q-pa-md">
            <div class="">
              <q-chip color="primary" text-color="white" class="q-mb-sm">
                {{ item.group }}
                <q-badge round color="secondary" text-color="primary" floating>{{ item.items.length }}</q-badge>
              </q-chip>

              <div v-for="(item) in item.items" :key="index" class="timeline-item">
                <div class="timeline-item-content">
                  <div class="timeline-item-content-title">{{ item.title }}</div>
                  <div class="timeline-item-content-timestamp">{{ item.timestamp }}</div>
                </div>
              </div>

            </div>
          </q-card>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

// Item is the type of each item in the timeline
type Item = {
  timestamp: string;
  title: string;
};

// Array of _Items to be displayed in the timeline
// this is a ref so it can be updated
let items: Array<Item> = [
  { timestamp: '2023-08-01T10:00:00Z', title: 'Event 1' },
  { timestamp: '2023-08-02T10:00:00Z', title: 'Event 1.2' },
  { timestamp: '2023-08-03T10:00:00Z', title: 'Event 1.3' },
  { timestamp: '2023-08-02T14:30:00Z', title: 'Event 2' },
  { timestamp: '2023-08-05T09:15:00Z', title: 'Event 3' },
  { timestamp: '2023-08-10T16:45:00Z', title: 'Event 4' },
  { timestamp: '2023-08-15T08:00:00Z', title: 'Event 5' },
  { timestamp: '2023-01-15T08:00:00Z', title: 'Event 6' },
  { timestamp: '2023-02-15T08:00:00Z', title: 'Event 6' },
  { timestamp: '2023-03-15T08:00:00Z', title: 'Event 6' },
]

const grouping = ref('day'); // Default grouping

type GroupedItems = {
  group: string
  items: Array<Item>
}


// group Item by the groupin value of day, week, month, quarter, year
// and return GroupedItems where day means all timestamps of that day, week means all timestamps of that week etc
const groupedItems = computed(() => {
  let myItems: Array<Item> = [];
  const groupedItems: Array<GroupedItems> = [];
  const groupedItemsMap = new Map<string, Array<Item>>();

  myItems = items;
  // sort myItems by timestamp
  myItems.sort((a, b) => {
    return new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime();
  });



  myItems.forEach((item) => {
    const groupKey = getGroupKey(item.timestamp);
    if (groupedItemsMap.has(groupKey)) {
      groupedItemsMap.get(groupKey)?.push(item);
    } else {
      groupedItemsMap.set(groupKey, [item]);
    }
  });

  groupedItemsMap.forEach((value, key) => {
    groupedItems.push({ group: key, items: value });
  });

  return groupedItems;
});



const setGrouping = (newGrouping) => {
  grouping.value = newGrouping;
};

const getGroupKey = (timestamp) => {
  const date = new Date(timestamp);
  switch (grouping.value) {
    case 'day':
      return date.toISOString().split('T')[0];
    case '3day':
      // split the month in as many 3 days peroid as possible and return year, month and period
      return date.getFullYear().toString() + '-' + (date.getMonth() + 1).toString() + '-' + Math.ceil(date.getDate() / 3).toString();
    case 'month':
      // return the month and year
      return date.toLocaleString('default', { month: 'long' }) + ' ' + date.getFullYear();
    case 'quarter':
      // return Q1 for the first 3 months of the year, Q2 for the next 3 months etc
      return 'Q' + Math.ceil((date.getMonth() + 1) / 3) + ' ' + date.getFullYear();
    case 'year':
      // return the year
      return date.getFullYear().toString();
    default:
      return date.toISOString().split('T')[0];
  }
};



</script>

<style>
/* Add your CSS styles here */
</style>
