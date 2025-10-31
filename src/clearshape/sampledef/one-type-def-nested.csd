type A = {
    s session?: string,
    t timestamp: enum {
        i isoTimestamp: [string, i64],
        u unixTimestampMillis: u64,
    },
    u user: {
        name: string,
        friends: {
            id: string,
        }[]
    }
}
