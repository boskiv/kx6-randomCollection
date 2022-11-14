import collection from 'k6/x/collection';
const collectionPath = 'collections';

export default function () {

    collection.createCollection('./collections')
    let item = collection.getRandomItem()
    console.log(JSON.stringify(item))

    let formData = collection.getRandomFormData()
    console.log(JSON.stringify(formData))
}
