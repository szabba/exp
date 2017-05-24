#lang rackjure


(require srfi/1)
(require "termprog.rkt")
(require "world-state.rkt")


(provide (contract-out [init (-> world-state? program?)]
                       [run-world! (-> world-state? void?)]))


;;;; Types
;; TODO: use non-dictionary for model


;; TODO: debug log
#;(struct debug-log
    (qualities covered-state)
    #:transparent)


(struct log-n-prompt
  (world log prompt scroll)
  #:transparent)


(struct done
  (last-state)
  #:transparent)


#;(type-alias Model
    (U LogNPrompt
       (Done Model)))


;;;; Boilerplate


#;(-> World Void)
(define (run-world! world)
  (run! (init world)))


#;(-> World (Program Model))
(define (init world)
  (program (init-model world) update view exit?))


#;(-> World Model)
(define (init-model world)
  (log-n-prompt world (list) (~> world world->prompt) 'prompt))

;;;; Helpers


#;(-> WorldState (Focus Option))
(define (world->prompt world)
  (match (~>> world world-state-available-options)
    [(list) (list)]
    [(list head rest ...)
     (focus (list) head rest)]))


#;(-> String x x)
(define (debug/log prefix v)
  (display prefix)
  (display ": ")
  (print v)
  (display "\r")
  (newline)
  v)


;;;; Update


#;(-> Key Model Model)
(define (update key model)
  (debug/log "update: key" key)
  (debug/log "update: model" model)

  (match (list model key)
    [(list (done _) _) model]

    [(list _ #\q) (done model)]

    [(list (log-n-prompt _ _ _ scroll-area) 'tab)
     (struct-copy log-n-prompt model
                  (scroll (~> scroll-area switch-scroll-area)))]

    [(list (log-n-prompt _ log _ 'log) 'up)
     (struct-copy log-n-prompt model
                  (log
                   (~>> log
                       (log-scroll focus-shift-up))))]

    [(list (log-n-prompt _ log _ 'log) 'down)
     (struct-copy log-n-prompt model
                  (log
                   (~>> log
                       (log-scroll focus-shift-down))))]

    [(list (log-n-prompt _ _ prompt 'prompt) 'up)
     (struct-copy log-n-prompt model
                  (prompt
                   (~>> prompt
                        (shift-options focus-shift-up))))]

    [(list (log-n-prompt _ _ prompt 'prompt) 'down)
     (struct-copy log-n-prompt model
                  (prompt
                   (~>> prompt
                        (shift-options focus-shift-down))))]

    [(list (log-n-prompt _ _ _ 'prompt) 'return)
     (~>> model
         log-focused-option
         (model-scroll-log focus-scroll-down)
         apply-focused-option
         fix-up-prompt)]

    [_ model]))


#;(-> (-> (Focus x) (Focus x)) LogNPrompt LogNPrompt)
(define (model-scroll-log f model)
  (struct-copy log-n-prompt model
               (log
                (~> model log-n-prompt-log f))))


#;(-> (-> (Focus x) (Focus x)) (U Null (Focus x)))
(define (log-scroll dir-f log)
  (match log
    [(list) log]

    [(focus _ _ _) (~> log dir-f)]))


#;(-> (U 'log 'prompt) (U 'log 'prompt))
(define (switch-scroll-area scroll-area)
  (match scroll-area
    ['log 'prompt]
    ['prompt 'log]))


#;(-> (-> (Focus x) (Focus x)) (U Null (Focus x)))
(define (shift-options dir opts)
  (match opts
    [(list) (list)]
    [(focus _ _ _) (dir opts)]))


#;(-> Model Model)
(define (log-focused-option model)
  (match model
    [(log-n-prompt _ _ (list) _) model]

    [(log-n-prompt _ log (focus _ opt _) _)

     (struct-copy log-n-prompt model
                  (log (~>> opt
                           (lines-for-option #:first-entry? (null? log))
                           (append-to-log log))))]))


#;(-> (U Null (Focus String)) (List String) (U Null (Focus String)))
(define (append-to-log log lines)

  (match (list log lines)
    [(list _ (list)) log]

    [(list (list) lines)
     (focus (drop-right lines 1) (last lines) (list))]

    [(list (focus _ _ _) lines)
     (focus-append-below log lines)]))


#;(-> Option #:first-entry? Boolean (List String))
(define (lines-for-option opt #:first-entry? [firts? #f])

  (define (indent line)
    (string-append "    " line))

  (match opt
    [(option title _ body-text _)

     (let [(lines (append (list title)
                          (list "")
                          (~>> body-text (map indent))))]

       (if first-entry?
           lines
           (append (list "") lines)))]))


#;(-> x Boolean)
(define first-entry? null?)


#;(-> LogNPrompt LogNPrompt)
(define (apply-focused-option model)
  (match model
    [(log-n-prompt _ _ (list) _) model]
    [(log-n-prompt world
     _ (focus _ opt _) _)
      (struct-copy log-n-prompt model
                  (world
                    (~>> opt
                        (world-state-apply-option world))))]))


#;(-> LogNPrompt LogNPrompt)
(define (fix-up-prompt model)
  (struct-copy log-n-prompt model
               (prompt
                (~>> model
                     log-n-prompt-world
                     world->prompt))))


;;;; View


#;(-> Int Int Model (List String))
(define (view w h model)
  (match model

    [(done last)
     (view w h last)]

    [(log-n-prompt _ log prompt _)
     (view/log-n-prompt w h log prompt)]))


#;(-> Int Int (U Null (Focus String)) (U Null (Focus Option)))
(define (view/log-n-prompt w h log prompt)

  (define bottom
    (list (make-string w #\-)
          (current-prompt-line prompt)))

  (define log-space
    (- h (length bottom)))

  (append (visible-log log-space log)
          bottom))


#;(-> Int (U Null (Focus String) (List String)))
(define (visible-log log-space log)
  (~> log
      log-above-prompt
      (prefix-with-upto log-space "")))


#;(-> (U Null (Focus String) (List String)))
(define (log-above-prompt log)
  (~> (match log
        [(list) (list)]
        [(focus above curr _) (cons curr above)])
      reverse))


#;(-> (List String) Int String (List String))
(define (prefix-with-upto l n p)
  (~> (append (make-list n p) l)
      (take-right-upto n)))


#;(-> (List String) Int (List String))
#;(define (drop-right-upto l n)
  (~>> n
       (min (length l))
       (drop-right l)))


#;(-> (List String) Int (List String))
(define (take-right-upto l n)
  (~>> n
       (min (length l))
       (take-right l)))


#;(-> (U Null (Focus Option) String))
(define (current-prompt-line prompt)
  (match prompt

    [(list) "! No actions available."]

    [(focus _ (option title _ _ _) _)
     (string-append "> " title)]))


#;(-> Model Boolean)
(define (exit? model)
  (match model
    [(done _) #t]
    [_ #f]))

;;;; Focus (a non-empty list zipper)


(struct focus (above on below) #:transparent)


#;(-> (Focus x) (Focus x))
(define (focus-shift-up foc)
  (match foc
    [(focus (list) _ _) foc]
    [(focus (list h t ...) m b)
     (focus t h (cons m b))]))


#;(-> (Focus x) (Focus x))
(define (focus-shift-down foc)
  (match foc
    [(focus _ _ (list)) foc]
    [(focus a m (list h t ...))
     (focus (cons m a) h t)]))


#;(-> (Focus x) (Focus x))
(define (focus-scroll-down foc)
  (match foc
    [(focus _ _ (list)) foc]
    [(focus a c b)
     (focus (append (drop-right b 1) (cons c a))
            (last b)
            (list))]))


#;(-> (Focus x) (Focus x))
#;(define (focus-append-above foc aboves)
  (match foc
    [(focus a c b)
     (focus (append (reverse aboves) a) c b)]))


#;(-> (Focus x) (Focus x))
(define (focus-append-below foc belows)
  (match foc
    [(focus a c b)
     (focus a c (append b belows))]))