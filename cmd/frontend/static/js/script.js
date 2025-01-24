document.addEventListener("DOMContentLoaded", (event) => {

    let fsCount = 0;
    let rowCount = 0;
    const newGroup = document.querySelector('#new_group');
    const clearAll = document.querySelector('#clear_all');
    const sendData = document.querySelector('#send');
    
    function addNewFieldset() {
        if (fsCount < 10) {
            const allFs = document.querySelector('#all_fieldsets');
            const fs = document.querySelector('.row-fieldset');
            const newFs = fs.cloneNode(true);
            newFs.style.display = 'block';
            newFs.setAttribute('data-fsid', fsCount);
            fsCount++;
            // new row groups with latitude & longitude            
            const rowGroups = newFs.querySelectorAll('.row-group');
            rowGroups.forEach(rgr => {
                rgr.setAttribute('data-rgid', rowCount);
                let siblingDist = rgr.nextElementSibling;
                siblingDist.setAttribute('data-dgid', rowCount);
                rowCount++;
            });

            // add row button
            const addRow = newFs.querySelector('.add-row');
            const rg = fs.querySelectorAll('.row-group')[0];
            const dg = fs.querySelectorAll('.dist-group')[0];
            addRow.addEventListener('click', function(e) {
                const newRg = rg.cloneNode(true);
                newRg.setAttribute('data-rgid', rowCount);
                const newDg = dg.cloneNode(true);
                newDg.setAttribute('data-dgid', rowCount);
                rowCount++;
                // remove row button
                const removeRow = newRg.querySelector('.remove-row');
                removeRow.addEventListener('click', function(e) {
                    e.preventDefault();
                    e.stopPropagation();
                    let parRow = e.target.parentElement;
                    let sDist = parRow.nextElementSibling;
                    parRow.remove();
                    sDist.remove();
                });
                // append divs with classes row-group & dist-group to fieldset
                newFs.append(newRg);
                newFs.append(newDg);
            });

            // remove row buttons
            const removeRows = newFs.querySelectorAll('.remove-row');
            removeRows.forEach(remRow => {
                remRow.addEventListener('click', function(e) {                    
                    let parentRow = remRow.parentElement;
                    let sibDist = parentRow.nextElementSibling;
                    parentRow.remove();
                    sibDist.remove();                    
                });                
            });

            // remove fieldset button
            const removeFs = newFs.querySelector('.remove-fs');
            removeFs.addEventListener('click', function(e) {
                const parentFs = e.target.parentElement.parentElement;
                parentFs.remove();
                fsCount--;
                if (fsCount < 10) {
                    if (newGroup.innerHTML.includes("max 10")) {
                        newGroup.innerHTML = 'New Group';
                    }
                }
                else {
                    if (!newGroup.innerHTML.includes("max 10")) {
                        newGroup.innerHTML += '<span style="color:red;"> (max 10)</span>';
                    }
                }
            });

            // append new fieldset        
            allFs.append(newFs);
            newGroup.innerHTML = 'New Group';
        }
        else {
            if (!newGroup.innerHTML.includes("max 10")) {
                newGroup.innerHTML += '<span style="color:red;"> (max 10)</span>';
            }
        }
    }
    
    addNewFieldset();


    newGroup.addEventListener('click', function(e) {
        e.preventDefault(); // Cancel the native event
        e.stopPropagation();// Don't bubble/capture the event any further
        addNewFieldset();
    });


    clearAll.addEventListener('click', function(e) {
        e.preventDefault();
        e.stopPropagation();
        const allFieldsets = document.querySelectorAll('fieldset');
        allFieldsets.forEach(fSet => {
            if (fSet.hasAttribute("data-fsid")) {
                fSet.remove();
            }
        });
        fsCount = 0;
        rowCount = 0;
        if (newGroup.innerHTML.includes("max 10")) {
            newGroup.innerHTML = 'New Group';
        }
        addNewFieldset();
    });


    sendData.addEventListener('click', function(e) {
        const parentForm = e.target.parentElement.parentElement;
        const fieldsets = parentForm.querySelectorAll('fieldset.row-fieldset');
        let allFieldsets = [];
        let fsObject = {}
        fieldsets.forEach(fdSet => {
            if (fdSet.checkVisibility()) allFieldsets.push(fdSet);
        });
        // create object of key - value: fieldset_id: array of rows with lat-long: {"0": ["0","1","2"]}
        allFieldsets.forEach(fldSet => {
            const fsId = fldSet.dataset.fsid;
            const rowGrps = fldSet.querySelectorAll('div.row-group');            
            if (rowGrps.length >= 2) {
                let rowIds = [];
                rowGrps.forEach(rg => {
                    const grId = rg.dataset.rgid;
                    rowIds.push(grId);
                });
                fsObject[fsId] = rowIds;
            }
        });
        // create pairs from arrays: {"0": [["0","1"], ["1","2"]]}
        pairObject = {}
        for (const [key, val] of Object.entries(fsObject)) {            
            let pair = val.reduce(function(result, value, index, array) {                
                result.push(array.slice(index, index + 2));
                return result;
            }, []);
            let pairs = [];
            pair.forEach(p => {
                if (p.length == 2) pairs.push(p)
            })
            pairObject[key] = pairs;
        }        
        // merge object values into array: [["0","1"], ["1","2"]]
        let allVals = [];        
        for (const val of Object.values(pairObject)) {
            allVals.push(val);
        }
        //const mergeVals = allVals.flat(1);
        const mergeVals = allVals.reduce(function(a, b){
            return a.concat(b);
        }, []);        
        // create object with key as row_id & value of array of coordinates: {"id": "0", "coords": {"point1": {"lat": 59.985363, "lng": 30.423156}, "point2": {"lat": 59.971793, "lng": 30.382472}}}
        let coordsGeo = {}
        let coordsObj = []    
        mergeVals.forEach(mv => {
            const row1 = document.querySelector(`[data-rgid="${mv[0]}"]`);
            const row2 = document.querySelector(`[data-rgid="${mv[1]}"]`);
            let cObj = {};
            let coords1 = {};
            let coords2 = {};
            const lat1 = parseFloat(row1.querySelector('.lat').value);
            const lng1 = parseFloat(row1.querySelector('.long').value);
            if (!isNaN(lat1) && !isNaN(lng1)) {
                coords1 = {"lat": lat1, "lng": lng1}
            }
            const lat2 = parseFloat(row2.querySelector('.lat').value);
            const lng2 = parseFloat(row2.querySelector('.long').value);
            if (!isNaN(lat2) && !isNaN(lng2)) {
                coords2 = {"lat": lat2, "lng": lng2}
            }            
            if (Object.keys(coords1).length && Object.keys(coords2).length) {                
                cObj["id"] = mv[0];
                cObj["coords"] = {
                    "point1": coords1,
                    "point2": coords2
                }
                coordsObj.push(cObj);
            }            
        });
        coordsGeo["locations"] = coordsObj
        
        // AJAX POST-data
        const headers = new Headers();
        headers.append("Content-Type", "application/json");

        const body = {
            method: "POST",
            body: JSON.stringify(coordsGeo),
            headers: headers,
        }

        fetch("/handle", body)
        .then((response) => response.json())
        .then((data) => {
            // clear all data-dgid
            const dElems = document.querySelectorAll('[data-dgid]');
            dElems.forEach(d => {
                const dDiv = d.querySelectorAll(".dist")[0];
                if (undefined != dDiv) {
                    dDiv.style.display = 'none';
                    dDiv.innerText = '';
                }                
            });
            // set distance to data-dgid
            const results = JSON.parse(data.result) // [{"id": "0", "distance": 7.050151045838122}]
            results.forEach(r => {
                const row = document.querySelector(`[data-dgid="${r.id}"]`);
                const distDiv = row.querySelectorAll(".dist")[0];
                if (undefined != distDiv) {
                    distDiv.style.display = 'block';
                    distDiv.innerText = `${r.distance} km`;
                }
            });
        })
        .catch((error) => {
            console.log(error)
        })
    })


})
