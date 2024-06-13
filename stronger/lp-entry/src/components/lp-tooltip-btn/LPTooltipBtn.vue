<template>
  <div class="tp-tooltip-btn">
    <v-tooltip 
      :bottom="bottom"
      :left="left"
      :right="right"
      :top="top"
    >
      <template v-slot:activator="{ on, attrs }">
        <v-btn
          :class="btnClass"
          :color="btnColor"
          :depressed="depressed"
          :dark="btnDark"
          :fab="fab"
          :icon="btnIcon"
          :outlined="outlined"
          :small="small"
          v-bind="attrs"
          v-on="on"
          @click="onClick"
        >
          <v-icon v-if="icon">{{ icon }}</v-icon>
          <slot name="text"></slot>
        </v-btn>
      </template>
      <!-- <slot name="tip"></slot> -->
      <span>{{ tip }}</span>
    </v-tooltip>
  </div>
</template>

<script>
  export default {
    name: 'LPToptipBtn',
    props: {
      bottom: {
        type: Boolean,
        default: false
      },

      btnClass: {
        type: String,
        required: false
      },

      btnColor: {
        type: String,
        required: false
      },

      btnDark: {
        type: Boolean,
        required: false
      },

      btnIcon: {
        type: Boolean,
        default: false
      },

      depressed: {
        type: Boolean,
        default: false
      },

      left: {
        type: Boolean,
        default: false
      },

      fab: {
        type: Boolean,
        default: false
      },

      icon: {
        type: String,
        required: false
      },

      outlined: {
        type: Boolean,
        default: false
      },

      right: {
        type: Boolean,
        default: false
      },

      small: {
        type: Boolean,
        default: false
      },

      tip: {
        type: String,
        required: true
      },

      top: {
        type: Boolean,
        default: false
      }
    },

    methods: {
      onClick(event) {
        this.$emit('click', event)
      }
    }
  }
</script>