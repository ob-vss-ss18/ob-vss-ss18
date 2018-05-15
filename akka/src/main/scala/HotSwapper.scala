import akka.actor.{
  ActorSystem, Actor, ActorRef,
  Props, PoisonPill
}
import language.postfixOps
import scala.concurrent.duration._

class HotSwapActor extends Actor {

  import context._

  def angry: Receive = {
    case "foo" =>
      println(s"${self.path} is already angry")
    case "bar" =>
      println(s"${self.path} becomes happy")
      become(happy)
  }

  def happy: Receive = {
    case "bar" =>
      println(s"${self.path} is already happy")
    case "foo" =>
      println(s"${self.path} becomes angry")
      become(angry)
  }

  def receive = {
    case "foo" =>
      println(s"${self.path} becomes angry for the first time")
      become(angry)
    case "bar" =>
      println(s"${self.path} becomes happy for the first time")
      become(happy)
  }
}

object HotSwapper extends App {
  val system = ActorSystem("hotswapper")

  val hotswapper = system.actorOf(Props[HotSwapActor], "hotswapper")

  import system.dispatcher

  system.scheduler.scheduleOnce(500 millis) {
    hotswapper ! "foo"
    hotswapper ! "foo"
    hotswapper ! "bar"
    hotswapper ! "bar"
    hotswapper ! "foo"
    hotswapper ! "bar"
    hotswapper ! PoisonPill
  }

}