package io.github.malczuuu.failbook

const val STATUS_UP = "UP"

data class Health(val status: String = STATUS_UP)
