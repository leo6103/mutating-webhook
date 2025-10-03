# Mutating Webhook

This webhook watches workspace pod creations and injects an init container that clones the current example repo and downloads the sync binary. Right now it always seeds `ray-example`, but we plan to look at workspace labels/metadata and pull the right repo dynamically once that information is carried through the CRDs.
