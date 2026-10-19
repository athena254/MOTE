# MOTE landing page

Live: https://mote-runtime.vercel.app

The public MOTE site is static HTML and CSS, with a small script for accessible
code tabs and copying. It has no dependencies, account forms, analytics, model
calls, serverless functions, or runtime secrets. Vercel may retain hosting logs.

## Local preview

From the repository root:

```sh
python -m http.server 8893 --bind 127.0.0.1 --directory site/public
```

Open http://127.0.0.1:8893. The code panels show integration examples; they do
not run an agent in the browser. The page links to MOTE's public GitHub
repository for source and downloads. The site serves only landing-page assets,
not repository documents or build output.

## Deployment

The Vercel project is `mote-runtime` in `kallmedis-projects`. Set the repository
root to `site`, the framework to Other, and the output directory to `public`.
Deployments use the CLI because Vercel rejected the repository connection during
setup. If you connect Git integration later, use `production-ready` as the
production branch. Link the project, then deploy from the repository root:

```sh
vercel deploy --prod --scope kallmedis-projects
```

Only `site/public` is served. The repository-root `.vercelignore` includes that
directory and `site/vercel.json`; inspect `vercel deploy --dry --json` before
publishing. Keep configuration, private documents, archives, credentials, and
build caches outside `public`. The CSP restricts scripts to the site itself,
blocks browser API network connections, and prevents other pages from framing it.

Before publishing, check for overflow on mobile and desktop. Test keyboard
navigation, code tabs, copying, FAQ expansion, and links. Check security headers,
the social image, and access to the production URL without authentication.
Update the version label for each MOTE release; site deployment is separate
from binary releases.
