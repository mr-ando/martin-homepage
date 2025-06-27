// @ts-check
import { defineConfig } from 'astro/config';

import tailwind from '@astrojs/tailwind';



import node from '@astrojs/node';
import react from '@astrojs/react';



// https://astro.build/config
export default defineConfig({
  output: 'static',
  integrations: [tailwind(), react(), node({mode: 'standalone'})],
  adapter: node({
    mode: 'standalone'
  }),
  outDir: 'build'
});