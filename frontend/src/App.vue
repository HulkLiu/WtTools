<script>
import { h } from "vue";
import { RouterLink } from "vue-router";
import { NIcon, NConfigProvider, NMenu, NTabs, NTabPane, NSpace, NLayout, NLayoutSider, NLayoutContent, NMessageProvider, NSwitch, NDivider } from 'naive-ui';

import {
  Home as HomeIcon,
  AlertCircleSharp as AboutIcon,
  SettingsSharp as SettingIcon,
  Barcode as OcrIcon,
} from '@vicons/ionicons5'

import { darkTheme } from 'naive-ui'
import mitt from './utils/event.js'

const renderIcon = (icon) => {
  return () => h(NIcon, null, { default: () => h(icon) });
}

export default {
  components: {
    NConfigProvider
  },
  data() {
    return {
      collapsed: false,
      switchTheme: false,
      myTheme: null,
      menuOptions: [
        {
          label: () => h(
            RouterLink,
            {
              to: {
                path: "/",
              }
            },
            { default: () => "主页" }
          ),
          key: "go-back-home",
          icon: renderIcon(HomeIcon)
        },
        {
          label: () => h(
            RouterLink,
            {
              to: {
                name: 'about',
                path: "/about",
              }
            },
            { default: () => "关于" }
          ),
          key: "go-back-about",
          icon: renderIcon(AboutIcon)
        },
        {
          label: () => h(
              RouterLink,
              {
                to: {
                  name: 'resource',
                  path: "/resource",
                }
              },
              { default: () => "资源管理" }
          ),
          key: "资源管理",
          icon: renderIcon(OcrIcon),
          path: "/resource",
        },
        {
          label: () => h(
              RouterLink,
              {
                to: {
                  name: 'task',
                  path: "/task",
                }
              },
              { default: () => "任务管理" }
          ),
          key: "任务管理",
          icon: renderIcon(OcrIcon),
          path: "/task",
        },
        {
          label: () => h(
              RouterLink,
              {
                to: {
                  name: 'setting',
                  path: "/setting",
                }
              },
              { default: () => "软件设置" }
          ),
          key: "软件设置",
          icon: renderIcon(SettingIcon),
          path: "/setting",

        },
      ],
      railStyle: ({
        focused,
        checked
      }) => {
        const style = {};
        if (checked) {
          style.background = "#4B9D5F";
          if (focused) {
            style.boxShadow = "0 0 0 2px #d0305040";
          }
        } else {
          style.background = "#000000";
          if (focused) {
            style.boxShadow = "0 0 0 2px #2080f040";
          }
        }
        return style;
      }
    }
  },
  mounted () {
    // 初始化项目时将主题保存在localStorage中
    localStorage.setItem('theme', 1)
  },
  methods: {
    changeTheme() {
      if (this.switchTheme) {
        this.myTheme = darkTheme
        localStorage.setItem("theme", 0)
        mitt.emit("theme","0")
      } else {
        this.myTheme = null
        localStorage.setItem("theme", 1)
        mitt.emit("theme","1")
      }
    }
  }
}

</script>

<template>
  <n-config-provider :theme="myTheme">
    <n-space vertical size="large">
      <n-layout has-sider position="absolute">
        <n-layout-sider bordered collapse-mode="width" :collapsed-width="80" :width="150" :collapsed="collapsed"
          show-trigger @collapse="collapsed = true" @expand="collapsed = false" style="--wails-draggable:drag; opacity: 1;">
            <n-menu :options="menuOptions" :collapsed-width="64" :collapsed-icon-size="22" style="margin-top: 40px;" />
            <div class="switchBtnPar">
              <n-divider />
              <n-switch :rail-style="railStyle" v-model:value="switchTheme" @update:value="changeTheme()"
                class="switchBtn">
                <template #checked>
                  亮
                </template>
                <template #unchecked>
                  暗
                </template>
              </n-switch>
            </div>
        </n-layout-sider>
        <n-layout-content>
          <n-message-provider>
            <router-view />
          </n-message-provider>
        </n-layout-content>

      </n-layout>
    </n-space>
  </n-config-provider>
</template>

<style>
.switchBtnPar {
  position: relative;
}

.switchBtn {
  position: absolute;
  left: 50%;
  transform: translate(-50%);
}
body {
  margin: 0;
}
</style>