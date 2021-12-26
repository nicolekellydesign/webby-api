ALTER TABLE gallery_items RENAME video_key TO embed_url;
UPDATE gallery_items SET embed_url = 'https://youtube.com/embed/' || embed_url WHERE embed_url IS NOT NULL;