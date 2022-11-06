import Client from "../client/client";
import User from "../models/user";
import Publication from "../models/publication";
import Source from "../models/source";

async function main() {
    let client = new Client('http://localhost:6060')

    let users: [User] = await client.users()
        .getUsers(0, 30)

    let sources: [Source] = await client.sources().getSources()
    if (!sources) {
        console.log('no sources available')
    }

    console.log('users:', users)

    for (let user of users) {
        try {
            await client.pubsub().subscribe(user.id, 'main')
        } catch (ex) {
            console.log(`failed to create subscription (${user.id}, main)`)
        }
    }

    console.log('here we go')

    await client.pubsub().publish(new Publication(
        sources[0].name,
        'greeting',
        null,
        ['main'],
        null,
        0
    ))
}


main()
    .then(() => {
    })
    .catch(err => {
        console.log('something went wrong:', err)
    })
