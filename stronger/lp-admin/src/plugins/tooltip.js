import TooltipComponent from '@/components/lp-tooltip/';

let instance;

const TooltipPlugin = {
  install(Vue) {
    // 如果已经存在实例，则直接返回
    if (instance) {
      return;
    }

    const TooltipConstructor = Vue.extend(TooltipComponent);

    instance = new TooltipConstructor();
    instance.$mount(document.createElement('div'));
    document.body.appendChild(instance.$el);

    // 添加全局方法
    Vue.prototype.$tooltipshow = (options) => {
      instance.show(options);
    };
    Vue.prototype.$tooltiphide = (options) => {
      instance.hide(options);
    };
  }
};

export default TooltipPlugin;