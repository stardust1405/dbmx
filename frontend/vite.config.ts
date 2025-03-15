import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [
		sveltekit(),
		{
			name: 'monaco-esm-fix',
			configureServer(server) {
			  server.middlewares.use((_req, res, next) => {
				res.setHeader('Cross-Origin-Embedder-Policy', 'require-corp');
				res.setHeader('Cross-Origin-Opener-Policy', 'same-origin');
				next();
			  });
			}
		  }
	],
	server: {
		host: true,
		port: 5173,
		strictPort: true,
	},
});
