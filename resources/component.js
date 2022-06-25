

document.addEventListener('click', function (e) {
    console.log(parseOutMethodAndParams(
        'Test(1,2,3, user.id)',
        {user: {id: 9}, $event: e}
    ));
})

function parseOutMethodAndParams (rawMethod, scope) {
    let method = rawMethod
    let params = []
    const parts = method.match(/(.*?)\((.*)\)/s)

    if (parts) {
        method = parts[1]

        let func = new Function(Object.keys(scope), `return (function () {
                for (var l=arguments.length, p=new Array(l), k=0; k<l; k++) {
                    p[k] = arguments[k];
                }
                return [].concat(p);
            })(${parts[2]})`)


        params = func(...Object.values(scope))
    }

    return {method, params}
}
