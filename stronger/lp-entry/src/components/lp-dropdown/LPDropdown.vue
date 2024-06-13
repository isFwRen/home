<template>
  <div class="lp-dropdown">
    <v-menu 
      :auto="auto"
      :bottom="bottom"
      :closeOnClick="closeOnClick"
      :closeOnContentClick="closeOnContentClick"
      :left="left"
      :offset-x="offsetX"
      :offset-y="offsetY"
      :position-x="positionX"
      :position-y="positionY"
      :right="right"
      :rounded="rounded"
      :tile="tile"
      :top="top"
      :z-index="zIndex"
      @input="onInput"
    >
      <template v-slot:activator="{ on, attrs }">
        <v-btn
          :absolute="absolute"
          :block="block"
          :bottom="bottom"
          :color="color"
          :depressed="depressed"
          :disabled="disabled"
          :fab="fab"
          :fixed="fixed"
          :icon="icon"
          :large="large"
          :left="left"
          :outlined="outlined"
          :plain="plain"
          :right="right"
          :rounded="rounded"
          :small="small"
          :tile="tile"
          :top="top"
          :x-large="larger"
          :x-small="smaller"
          v-bind="attrs"
          v-on="on"
        >
          <slot></slot>
        </v-btn>
      </template>
      <v-list>
        <v-list-item
          v-for="(item, index) in items"
          :key="`z_dropdown_${ index }`"
          link
          @click="onClick($event, item)"
        >
          <v-list-item-title>{{ item.label }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
  </div>
</template>

<script>
  import { tools } from 'vue-rocket'

  export default {
    name: 'LPDropdown',

    props: {
      auto: {
        type: Boolean,
        default: false
      },

      closeOnClick: {
        type: Boolean,
        default: true
      },

      closeOnContentClick: {
        type: Boolean,
        default: true
      },

      offsetX: {
        type: Boolean,
        default: false
      },

      offsetY: {
        type: Boolean,
        default: false
      },

      options: {
        type: Array,
        required: false
      },

      positionX: {
        type: Number,
        required: false
      },

      positionY: {
        type: Number,
        required: false
      },

      top: {
        type: Boolean,
        default: false
      },  

      zIndex: {
        type: [Number, String],
        required: false
      },

      // button
      absolute: {
        type: Boolean,
        default: false
      },

      block: {
        type: Boolean,
        default: false
      },

      bottom: {
        type: Boolean,
        default: false
      },

      color: {
        type: String,
        required: false
      },

      className: {
        type: String,
        required: false
      },

      dark: {
        type: Boolean,
        default: false
      },

      depressed: {
        type: Boolean,
        default: false
      },

      disabled: {
        type: Boolean,
        default: false
      },

      elevation: {
        type: [Number, String],
        required: false
      },

      fab: {
        type: Boolean,
        default: false
      },

      fixed: {
        type: Boolean,
        default: false
      },

      height: {
        type: [Number, String],
        required: false
      },

      icon: {
        type: Boolean,
        default: false
      },

      large: {
        type: Boolean,
        default: false
      },

      larger: {
        type: Boolean,
        default: false
      },

      left: {
        type: Boolean,
        default: false
      },

      loading: {
        type: Boolean,
        default: false
      },

      outlined: {
        type: Boolean,
        default: false
      },

      plain: {
        type: Boolean,
        default: false
      },

      right: {
        type: Boolean,
        default: false
      },

      rounded: {
        type: Boolean,
        default: false
      },

      small: {
        type: Boolean,
        default: false
      },

      smaller: {
        type: Boolean,
        default: false
      },

      text: {
        type: Boolean,
        default: false
      },

      tile: {
        type: Boolean,
        default: false
      },

      width: {
        type: [Number, String],
        required: false
      },

      unlocked: {
        type: Boolean,
        default: false
      }
    },

    data() {
      return {
        items: []
      }
    },

    methods: {
      onClick(event, { value }) {
        event.customValue = value
        this.$emit('click', event)
      },

      onInput(input) {
        this.$emit('input', input)
      },

      _setOptions() {
        if(tools.isArray(this.options) && tools.isYummy(this.options)) {
          this.items = this.options
        }
      }
    },

    watch: {
      options: {
        handler() {
          this._setOptions()
        },
        immediate: true
      }
    }
  }
</script>