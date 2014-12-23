package main

type Storage interface {
    Get(query AeroGet) *AeroResponse
    BatchGet(query AeroGet) *[]*AeroResponse
    Put(new_item AeroPut) bool
    Delete(item AeroDelete) bool
    Query(entry AeroQuery) *[]*AeroResponse
    CreateIndex(index AeroIndex) bool
    DropIndex(index AeroIndex) bool
}
