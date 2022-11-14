import collection from 'k6/x/collection';
const collectionPath = 'collections';

export default function () {

    collection.createCollection('./collections')
    let data = collection.getRandomItem()
    console.log(JSON.stringify(data))
}
