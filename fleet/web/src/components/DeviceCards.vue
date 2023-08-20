<script setup>
import DeviceStateBadge from './DeviceStateBadge.vue';
import DeviceStatusBadge from './DeviceStatusBadge.vue';

defineProps(['title'])
</script>

<template>
  <ul class="cards">
    <li>
      <a href="" class="card">
        <img src="img/vendor/vendor_huawei.svg" class="card__image" alt="" />

        <div class="card__icons">
          <DeviceStateBadge state="ACTIVE" />
          <DeviceStatusBadge status="REACHABLE" />
        </div>
        <div class="card__overlay">
          <div class="card__header">

            <q-btn round size="lg" color="primary"><q-icon name="fa-solid fa-circle-nodes" /></q-btn>
            <div class="card__header-text">
              <h3 class="card__title">sait-hea-a3</h3>
              <span class="card__tagline">10.5.5.23</span>
              <span class="card__status">1 hour ago</span>
            </div>
          </div>
          <div class="card__description">
            <q-list dense>
              <q-item>
                <q-item-section avatar>
                  <q-icon name="public" size="sm"/>
                </q-item-section>

                <q-item-section>
                  <q-item-label class="heading">Management IP</q-item-label>
                  <q-item-label caption>
                    10.5.5.23
                  </q-item-label>
                </q-item-section>

              </q-item>

              <q-item>
                <q-item-section avatar>
                  <q-icon name="hub" />
                </q-item-section>
                <q-item-section>
                  <q-item-label class="heading">Network Region</q-item-label>
                  <q-item-label caption>
                    VX_SAIT1
                  </q-item-label>
                </q-item-section>

              </q-item>


            </q-list>
          </div>
        </div>
      </a>
    </li>
    <li>
      <a href="" class="card">

        <img src="https://i.imgur.com/2DhmtJ4.jpg" class="card__image" alt="" />
        <div class="card__icons">
          <DeviceStateBadge state="ACTIVE" />
          <DeviceStatusBadge status="REACHABLE" />
        </div>
        <div class="card__overlay">

          <div class="card__header">
            <svg class="card__arc" xmlns="http://www.w3.org/2000/svg">
              <path />
            </svg>
            <img class="card__thumb" src="https://i.imgur.com/sjLMNDM.png" alt="" />
            <div class="card__header-text">
              <h3 class="card__title">kim Cattrall</h3>
              <span class="card__status">3 hours ago</span>
            </div>

          </div>
          <p class="card__description">Lorem ipsum dolor sit amet consectetur adipisicing elit. Asperiores, blanditiis?
          </p>
        </div>
      </a>
    </li>
    <li>
      <a href="" class="card">
        <img src="https://i.imgur.com/oYiTqum.jpg" class="card__image" alt="" />
        <div class="card__overlay">
          <div class="card__header">
            <svg class="card__arc" xmlns="http://www.w3.org/2000/svg">
              <path />
            </svg>
            <img class="card__thumb" src="https://i.imgur.com/7D7I6dI.png" alt="" />
            <div class="card__header-text">
              <h3 class="card__title">Jessica Parker</h3>
              <span class="card__tagline">Lorem ipsum dolor sit amet consectetur</span>
              <span class="card__status">1 hour ago</span>
            </div>
          </div>
          <p class="card__description">Lorem ipsum dolor sit amet consectetur adipisicing elit. Asperiores, blanditiis?
          </p>
        </div>
      </a>
    </li>
    <li>
      <a href="" class="card">
        <img src="https://i.imgur.com/2DhmtJ4.jpg" class="card__image" alt="" />
        <div class="card__overlay">
          <div class="card__header">
            <svg class="card__arc" xmlns="http://www.w3.org/2000/svg">
              <path />
            </svg>
            <img src="img/vendor/vendor_huawei.svg" class="card__thumb" alt="" />
            <div class="overlay"></div>
            <div class="card__header-text">
              <h3 class="card__title">kim Cattrall</h3>
              <span class="card__status">3 hours ago</span>
            </div>
          </div>
          <p class="card__description">Lorem ipsum dolor sit amet consectetur adipisicing elit. Asperiores, blanditiis?
          </p>
        </div>
      </a>
    </li>
  </ul>
</template>

<style lang="sass">
:root
  --surface-color: #fff
  --curve: 40


.cards
  display: grid
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr))
  gap: 2rem
  margin: 4rem 5vw
  padding: 0
  list-style-type: none



.card
  position: relative
  display: block
  height: 100%
  background-color: var(--surface-color)
  box-shadow: 0 10px 30px -10px rgba(0, 0, 0, 0.2)
  border-radius: calc(var(--curve) * 1px)
  overflow: hidden
  text-decoration: none
  padding:15px


.card__image
  z-index: 10
  width: 100%
  height: auto
  border-radius: calc(var(--curve) * 1px)


.card__icons
  margin: 15px
  position: absolute
  right: 0px
  top: 0px
  z-index: 2

.card__overlay
  position: absolute
  bottom: 0
  left: 0
  right: 0
  z-index: 1
  border-radius: calc(var(--curve) * 1px)
  background-color: var(--surface-color)
  transform: translateY(100%)
  transition: .2s ease-in-out


.card:hover .card__overlay
  transform: translateY(0)


.card__header
  position: relative
  display: flex
  align-items: center
  gap: 2em
  padding: 2em
  border-radius: calc(var(--curve) * 1px) 0 0 0
  background-color: var(--surface-color)
  transform: translateY(-100%)
  transition: .2s ease-in-out

.card__header-text h3
  line-height: 1

.card__arc
  width: 80px
  height: 80px
  position: absolute
  bottom: 100%
  right: 0
  z-index: 1


.card__arc path
  fill: var(--surface-color)
  d: path("M 40 80 c 22 0 40 -22 40 -40 v 40 Z")


.card:hover .card__header
  transform: translateY(0)


.card__thumb
  flex-shrink: 0
  width: 50px
  height: 50px
  border-radius: 00%



.card__title
  font-size: 1em
  margin: 0 0 .3em
  color: #6A515E


.card__tagline
  display: block
  margin: 0.3em 0
  font-family: $typography-font-family
  font-size: .8em
  color: rgb(72, 72, 72)


.card__status
  font-size: .9em
  color: $primary


.card__description
  padding: 0 2em 2em
  margin: 0
  color: #000000
  display: -webkit-box
  -webkit-box-orient: vertical
  -webkit-line-clamp: 3
  overflow: hidden

</style>
