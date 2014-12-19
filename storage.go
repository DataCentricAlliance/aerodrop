package main

type Storage interface {
    Get(query AeroPK) *AeroResponse
    BatchGet(query AeroPK) *[]*AeroResponse
    Put(new_item AeroNew) bool
    Delete(item AeroDelete) bool
    Query(entry AeroQuery) *[]*AeroResponse
    CreateIndex(index AeroIndex) bool
    DropIndex(index AeroIndex) bool
}
