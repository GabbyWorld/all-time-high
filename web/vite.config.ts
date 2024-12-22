import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { nodePolyfills } from "vite-plugin-node-polyfills";
import postcssPxtorem from "postcss-pxtorem";
import path from "path";
import mkcert from "vite-plugin-mkcert"
import { createFilter } from '@rollup/pluginutils';

const cacheBusterPlugin = () => {
  const filter = createFilter(['**/*.tsx', '**/*.ts', '**/*.js', '**/*.jsx']);

  return {
    name: 'cache-buster-plugin',
    transform(code: string, id: string) {
      if (!filter(id)) return;

      // Replace resource URLs with cache-busted URLs
      const result = code.replace(/(src|href)="(https:\/\/[^"]+)"/g, (match, p1, p2) => {
        const cacheBustedUrl = `${p2}?v=${Date.now()}`;
        return `${p1}="${cacheBustedUrl}"`;
      });

      return {
        code: result,
        map: null, // If you have source maps, add them here
      };
    },
  };
};


export default defineConfig(config => {
  return {
    plugins: [
      react(), 
      nodePolyfills(),
      mkcert(),
      // cacheBusterPlugin()
    ],
    server: {
      port: 3000,
      host:"0.0.0.0",
      https: true
    },
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "src"),
      },
    },
    css: {
      // postcss: {
      //   plugins: [
          
      //     postcssPxtorem({
      //       rootValue: 16, // 根据设计稿设置根元素字体大小
      //       propList: ["*"], // 可以设置需要转换的属性列表，['*'] 表示所有属性
      //       unitPrecision: 5, // 转换后的小数精度
      //       selectorBlackList: [], // 可以设置不需要转换的选择器
      //       replace: true, // 替换而不是添加回退
      //       mediaQuery: false, // 允许在媒体查询中转换
      //       minPixelValue: 0, // 设置最小需要转换的数值
      //     }),
      //   ],
      // },
    },
    base: ((process.env.GITHUB_REPOSITORY ?? "") + "/").match(/(\/.*)/)?.[1],
    }
});

