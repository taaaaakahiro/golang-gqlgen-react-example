async function fetchGraphQL(gqlQuery, params) {
    const response = await fetch('http://0.0.0.0:8080/query', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Origin': 'http://0.0.0.0:3001',
        },
        mode: 'cors',
        cache: 'no-cache',
        body: JSON.stringify({
            query: gqlQuery,
            variables: params
        })
    });
    return await response.json();
}

export default fetchGraphQL;