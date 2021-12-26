ALTER TABLE gallery_items RENAME embed_url TO video_key;
UPDATE gallery_items SET video_key = LTRIM(video_key, 'https://youtube.com/embed/');