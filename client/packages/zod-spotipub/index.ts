import { z } from 'zod';

const entity = z.object({
    id: z.string(),
    name: z.string(),
    href: z.string(),
    uri: z.string(),
    external_url: z.object({
        spotify: z.string()
    })
});

export const payloadSchema = z.object({
    timestamp: z.number(),
    progress_ms: z.number(),
    is_playing: z.boolean(),
    item: entity.extend({
        artists: entity.array(),
        duration_ms: z.number(),
        explicit: z.boolean(),
        type: z.literal("track"),
        popularity: z.number(),
        album: entity.extend({
            artists: entity.array(),
            album_type: z.string(),
            images: z.object({
                height: z.number(),
                width: z.number(),
                url: z.string()
            }).array()
        }),
    }).nullable()
});

export type Payload = z.infer<typeof payloadSchema>;