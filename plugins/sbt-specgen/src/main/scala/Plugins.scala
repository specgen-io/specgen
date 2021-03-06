package spec

import sbt._
import sbt.Keys._
import Specgen._

object SpecKeys {
  lazy val specFile = settingKey[File]("Path to service specification file")
  lazy val specSwagger = settingKey[File]("Path to folder where swagger specification should be generated")
  lazy val specGeneratePath = settingKey[File]("Path to generate source code into")
  lazy val specServicesPath = settingKey[File]("Path to scaffold services code")

  lazy val specgenServicePlay = taskKey[Seq[File]]("Run service Play code generation for spec")
  lazy val specgenModels = taskKey[Seq[File]]("Run service models code generation for spec")
  lazy val specgenClientSttp = taskKey[Seq[File]]("Run client Sttp code generation for spec")

  lazy val specModelDependencies = Seq(
    "io.circe" %% "circe-core" % "0.12.3",
    "io.circe" %% "circe-generic-extras" % "0.12.2",
    "io.circe" %% "circe-parser" % "0.12.3",
    "com.beachape" %% "enumeratum" % "1.5.13",
    "com.beachape" %% "enumeratum-circe" % "1.5.22",
  )

  lazy val specPlayDependencies = Seq(
    "io.circe" %% "circe-core" % "0.12.3",
    "io.circe" %% "circe-generic-extras" % "0.12.2",
    "io.circe" %% "circe-parser" % "0.12.3",
    "com.beachape" %% "enumeratum" % "1.5.13",
    "com.beachape" %% "enumeratum-circe" % "1.5.22",

    "com.typesafe.play" %% "play" % "2.8.1",
    "org.webjars" % "swagger-ui" % "3.22.2",
  )

  lazy val specSttpDependencies = Seq(
    "io.circe" %% "circe-core" % "0.12.3",
    "io.circe" %% "circe-generic-extras" % "0.12.2",
    "io.circe" %% "circe-parser" % "0.12.3",
    "com.beachape" %% "enumeratum" % "1.5.13",
    "com.beachape" %% "enumeratum-circe" % "1.5.22",

    "org.slf4j" % "slf4j-api" % "1.7.28",
    "com.softwaremill.sttp" %% "core" % "1.7.1",
  )
}

object CommonKeys extends AutoPlugin {
  val autoImport = SpecKeys
}

object SpecPlay extends AutoPlugin {
  import SpecKeys._

  private def specgenTask = Def.task {
    serviceScalaPlay(
      sLog.value,
      specFile.value,
      specSwagger.value,
      specGeneratePath.value,
      specServicesPath.value,
    )
  }

  override val projectSettings = Seq(
    specFile := file("spec.yaml"),
    specSwagger := baseDirectory.value / "public",
    specGeneratePath := (sourceManaged in Compile).value / "spec",
    specServicesPath := (scalaSource in Compile).value / "services",
    specgenServicePlay := specgenTask.value,
    sourceGenerators in Compile += specgenServicePlay,
    mappings in (Compile, packageSrc) ++= {(specgenServicePlay in Compile) map { sourceFiles =>
      sourceFiles map { sourceFile => (sourceFile, sourceFile.getName)}
    }}.value
  )
}

object SpecCirce extends AutoPlugin {
  import SpecKeys._

  private def specgenTask = Def.task {
    modelsScalaCirce(
      sLog.value,
      specFile.value,
      specGeneratePath.value,
    )
  }

  override val projectSettings = Seq(
    specFile := file("spec.yaml"),
    specGeneratePath := (sourceManaged in Compile).value / "spec",
    specgenModels := specgenTask.value,
    sourceGenerators in Compile += specgenModels,
    mappings in (Compile, packageSrc) ++= {(specgenModels in Compile) map { sourceFiles =>
      sourceFiles map { sourceFile => (sourceFile, sourceFile.getName)}
    }}.value
  )
}

object SpecSttp extends AutoPlugin {
  import SpecKeys._

  private def specgenTask = Def.task {
    clientSttp(
      sLog.value,
      specFile.value,
      specGeneratePath.value,
    )
  }

  override val projectSettings = Seq(
    specFile := file("spec.yaml"),
    specGeneratePath := (sourceManaged in Compile).value / "spec",
    specgenClientSttp := specgenTask.value,
    sourceGenerators in Compile += specgenClientSttp,
    mappings in (Compile, packageSrc) ++= {(specgenClientSttp in Compile) map { sourceFiles =>
      sourceFiles map { sourceFile => (sourceFile, sourceFile.getName)}
    }}.value
  )
}