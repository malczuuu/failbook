import com.diffplug.spotless.LineEnding
import com.diffplug.spotless.kotlin.KtfmtStep.TrailingCommaManagementStrategy

plugins {
  id("com.gradleup.shadow").version("9.2.2")
  id("com.diffplug.spotless").version("8.0.0")
  kotlin("jvm").version("2.2.20")
}

group = "io.github.malczuuu.failbook"

version = "1.0.0-SNAPSHOT"

repositories { mavenCentral() }

dependencies {
  implementation("io.javalin:javalin:6.7.0")
  implementation("io.javalin:javalin-micrometer:6.7.0")

  implementation("io.micrometer:micrometer-registry-prometheus:1.15.5")

  implementation("org.slf4j:slf4j-api:2.0.17")
  implementation("ch.qos.logback:logback-classic:1.5.20")
  implementation("co.elastic.logging:logback-ecs-encoder:1.7.0")

  implementation("com.fasterxml.jackson.core:jackson-databind:2.20.0")

  testImplementation(kotlin("test"))
}

spotless {
  format("misc") {
    target("**/.dockerignore", "**/.gitattributes", "**/.gitignore")

    trimTrailingWhitespace()
    leadingTabsToSpaces(4)
    endWithNewline()
    lineEndings = LineEnding.UNIX
  }

  kotlin {
    target("**/*.kt", "**/*.kts")

    ktfmt("0.59").metaStyle().configure {
      it.setRemoveUnusedImports(true)
      it.setTrailingCommaManagementStrategy(TrailingCommaManagementStrategy.NONE)
    }
    endWithNewline()
    lineEndings = LineEnding.UNIX
  }
}

tasks.jar { enabled = false }

tasks.shadowJar {
  manifest { attributes["Main-Class"] = "io.github.malczuuu.failbook.MainKt" }
  archiveClassifier = null
}

tasks.withType<Test>().configureEach { useJUnitPlatform() }
