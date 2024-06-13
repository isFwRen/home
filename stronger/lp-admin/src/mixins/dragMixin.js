
export default {
  data() {
    return {
      dragBox: null,
      isStart: false,
      position: {
        x: 0,
        y: 0
      },
      startPos: {
        x: 0,
        y: 0
      }
    }
  },
  computed: {
    moveStyle() {
      return `transform: translate(${this.position.x}px,${this.position.y}px);touch-action:none;`;
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.bindEventListenr(this.$refs.dragTarget, {});
      this.dragBox = this.$refs.dragEle;
    });
  },
  methods: {
    handleEvent(event) {
      event.stopPropagation();
      event.preventDefault();
    },
    start(event) {
      this.startPos.x = event.clientX - this.position.x;
      this.startPos.y = event.clientY - this.position.y;
      this.isStart = true;
      this.handleEvent(event);
    },
    move(event) {
      if (!this.isStart) return;
      let { x, y } = this.position;
      x = event.clientX - this.startPos.x;
      y = event.clientY - this.startPos.y;
      this.position.x = x;
      this.position.y = y;
      this.handleEvent(event);
    },
    up(event) {
      this.isStart = false;
      this.handleEvent(event);
    },
    useEventListener(target, type, handler, config) {
      if (!target) {
        return;
      }
      target.addEventListener(type, handler, config);
    },
    bindEventListenr(ele, options) {
      const config = { capture: options.capture ?? true };
      this.useEventListener(ele, "pointerdown", this.start, config);
      this.useEventListener(window, "pointermove", this.move, config);
      this.useEventListener(window, "pointerup", this.up, config);
    },
    unbindEventListener(ele, options) {
      const config = { capture: options.capture ?? true };
      this.removeEventListener(ele, "pointerdown", this.start, config);
      this.removeEventListener(window, "pointermove", this.move, config);
      this.removeEventListener(window, "pointerup", this.up, config);
    },
    removeEventListener(target, type, handler, config) {
      if (!target) {
        return;
      }
      target.removeEventListener(type, handler, config);
    },
  }
}