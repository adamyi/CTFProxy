import json
import requests
import os

from flask import request

from CTFd.cache import clear_standings
from CTFd.plugins import register_plugin_assets_directory, register_plugin_script
from CTFd.plugins.challenges import get_chal_class
from CTFd.models import db, Challenges, Flags, Pages, Users, Solves
from CTFd.utils import user as current_user
from CTFd.utils.dates import ctf_paused, ctftime
from CTFd.utils.user import authed, get_current_team, get_current_user, is_admin
from CTFd.utils import config
from CTFd.utils.logging import log

FLAGANIZER_DESCRIPTION_PREFIX = "[DO NOT MODIFY AND KEEP HIDDEN] Flaganizer controlled challenge. "
FLAGANIZER_VERIFY_ENDPOINT = "https://flaganizer." + os.environ[
    "CTF_DOMAIN"] + "/verify"


def load(app):
  # set_index_content()
  # import_users()
  register_plugin_assets_directory(app, base_path='/plugins/flaganizer/assets/')
  register_plugin_script('/plugins/flaganizer/assets/flaganizer.js')

  @app.route('/flaganizer-submit', methods=['POST'])
  def flaganizer_submit():
    if authed() is False:
      return {
          "success": True,
          "data": {
              "status": "authentication_required"
          }
      }, 403

    if request.content_type != "application/json":
      request_data = request.form
    else:
      request_data = request.get_json()

    if ctf_paused():
      return (
          {
              "success": True,
              "data": {
                  "status": "paused",
                  "message": "{} is paused".format(config.ctf_name()),
              },
          },
          403,
      )

    user = get_current_user()
    team = get_current_team()

    # TODO: Convert this into a re-useable decorator
    if config.is_teams_mode() and team is None:
      abort(403)

    kpm = current_user.get_wrong_submissions_per_minute(user.account_id)

    frsp = requests.post(
        FLAGANIZER_VERIFY_ENDPOINT,
        data={
            "flag": request_data.get("submission", "")
        },
        headers={
            "X-CTFProxy-SubAcc-JWT": request.headers.get("X-CTFProxy-JWT")
        }).json()
    if frsp["Success"] == 0:
      if ctftime() or current_user.is_admin():
        placeholder_challenge = Challenges.query.filter_by(
            name="wrong submission").first()
        if placeholder_challenge is None:
          placeholder_challenge = Challenges(
              name="wrong submission",
              description=FLAGANIZER_DESCRIPTION_PREFIX +
              "a placeholder challenge for unrecognized flags",
              value=0,
              category="misc",
              state="hidden",
              max_attempts=0)
          db.session.add(placeholder_challenge)
          db.session.commit()
          db.session.close()
          placeholder_challenge = Challenges.query.filter_by(
              name="wrong submission").first()
        chal_class = get_chal_class(placeholder_challenge.type)
        if placeholder_challenge is not None:
          chal_class.fail(
              user=user,
              team=team,
              challenge=placeholder_challenge,
              request=request)
          clear_standings()
      log(
          "submissions",
          "[{date}] {name} submitted {submission} via flaganizer with kpm {kpm} [WRONG]",
          submission=request_data.get("submission", "").encode("utf-8"),
          kpm=kpm,
      )
      return {
          "success": True,
          "data": {
              "status": "incorrect",
              "message": frsp["Message"]
          },
      }

    challenge = Challenges.query.filter_by(
        description=FLAGANIZER_DESCRIPTION_PREFIX + frsp["Flag"]["Id"]).first()
    if challenge is None:
      challenge = Challenges(
          name=frsp["Flag"]["DisplayName"],
          description=FLAGANIZER_DESCRIPTION_PREFIX + frsp["Flag"]["Id"],
          value=frsp["Flag"]["Points"],
          category=frsp["Flag"]["Category"],
          state="hidden",
          max_attempts=0)

      db.session.add(challenge)
      db.session.commit()
    challenge_id = challenge.id

    if challenge.state == "locked":
      db.session.close()
      abort(403)

    if challenge.requirements:
      requirements = challenge.requirements.get("prerequisites", [])
      solve_ids = (
          Solves.query.with_entities(Solves.challenge_id).filter_by(
              account_id=user.account_id).order_by(
                  Solves.challenge_id.asc()).all())
      solve_ids = set([solve_id for solve_id, in solve_ids])
      prereqs = set(requirements)
      if solve_ids >= prereqs:
        pass
      else:
        db.session.close()
        abort(403)

    chal_class = get_chal_class(challenge.type)

    if kpm > 10:
      if ctftime():
        chal_class.fail(
            user=user, team=team, challenge=challenge, request=request)
      log(
          "submissions",
          "[{date}] {name} submitted {submission} on {challenge_id} with kpm {kpm} [TOO FAST]",
          submission=request_data.get("submission", "").encode("utf-8"),
          challenge_id=challenge_id,
          kpm=kpm,
      )
      # Submitting too fast
      db.session.close()
      return (
          {
              "success": True,
              "data": {
                  "status": "ratelimited",
                  "message": "You're submitting flags too fast. Slow down.",
              },
          },
          429,
      )

    solves = Solves.query.filter_by(
        account_id=user.account_id, challenge_id=challenge_id).first()

    # Challenge not solved yet
    if not solves:
      status, message = chal_class.attempt(challenge, request)
      chal_class.solve(
          user=user, team=team, challenge=challenge, request=request)
      clear_standings()
      log(
          "submissions",
          "[{date}] {name} submitted {submission} on {challenge_id} via flaganizer with kpm {kpm} [CORRECT]",
          submission=request_data.get("submission", "").encode("utf-8"),
          challenge_id=challenge_id,
          kpm=kpm,
      )
      db.session.close()
      return {
          "success": True,
          "data": {
              "status": "correct",
              "message": "Successfully submitted!"
          },
      }
    else:
      log(
          "submissions",
          "[{date}] {name} submitted {submission} on {challenge_id} via flaganizer with kpm {kpm} [ALREADY SOLVED]",
          submission=request_data.get("submission", "").encode("utf-8"),
          challenge_id=challenge_id,
          kpm=kpm,
      )
      db.session.close()
      return {
          "success": True,
          "data": {
              "status": "already_solved",
              "message": "You already solved this",
          },
      }
