#lang rackjure

(require srfi/1)
(require "termprog.rkt")
(require "world-state.rkt")

(provide (contract-out [init (-> world-state? program?)]
                       [run-world! (-> world-state? void?)]))


(define (run-world! world)
  (run! (init world)))

(define (init world)
  (program (init-model world) update view exit?))

(define (init-model world)
  (~> {'world world 'scroll 'prompt 'options (~>> world world->options-focus) 'log (list) 'upscroll 0 'done #f}))

(define (world->options-focus world)
  (match (~>> world world-state-available-options)
    [(list) (list)]
    [(list head rest ...)
     (focus (list) head rest)]))


(define (update key model)
  (match (list (model 'scroll) key)
    [(list _ 'escape) (model 'done #t)]
    [(list _ #\q) (model 'done #t)]

    [(list _ 'pgup) (~> model switch-scroll-area)]
    [(list _ 'pgdn) (~> model switch-scroll-area)]

    [(list 'prompt 'up)

     (~>> model
          'options
          (shift-options focus-shift-up)
          (dict-set model 'options))]

    [(list 'prompt 'down)

     (~>> model
          'options
          (shift-options focus-shift-down)
          (model 'options))]
    
    [(list 'prompt 'return)
     (~> model
         log-focused-option
         apply-focused-option
         fix-up-options)]

    [(list 'log 'up)
     
     (~>> model
          (add-to-upscroll 1)
          fix-upscroll-range)]

    [(list 'log 'down)
     (~>> model
          (add-to-upscroll -1)
          fix-upscroll-range)]
    
    [_ model]))

(define (add-to-upscroll n model)
  (~>> model
       'upscroll
       (+ n)
       (model 'upscroll)))

(define (fix-upscroll-range model)
  (define upscroll (model 'upscroll))
  (define log-len (~>> model 'log length))
  
  (~>> (cond ((< upscroll 0) 0)
             ((<= log-len upscroll) (- log-len 1))
             (else upscroll))
       (model 'upscroll)))

(define (switch-scroll-area model)
  (model 'scroll
         (match (model 'scroll)
           ['log 'prompt]
           ['prompt 'log])))

(define (shift-options dir opts)
  (match opts
    [(list) (list)]
    [(focus _ _ _) (dir opts)]))

(define (log-focused-option model)
  
  (match (model 'options)
    [(list) model]

    [(focus _ opt _)
     (let [(log (model 'log))]

       (~>> (lines-for-option opt (first-entry? log))
            (append log)
            (model 'log)))]))

(define (lines-for-option opt first-entry?)
  (match opt
    [(option title _ body-text _)

     (let [(lines (append (list title)
                          (list "")
                          (~>> body-text
                               
                               (map (λ (line) (string-append "    " line))))))]

       (if first-entry?
           lines
           (append (list "") lines)))]))

(define first-entry? null?)

(define (apply-focused-option model)
  
  (let [(world (model 'world))
        (opts (model 'options))]
    
    (match opts
      [(list) model]
      [(focus _ opt _)
       (~>> opt
            (world-state-apply-option world)
            (model 'world))])))

(define (fix-up-options model)
  (~>> (model 'world)
       world->options-focus
       (model 'options)))


(define (view w h model)

  (define log (model 'log))
  (define log-len (length log))

  (define bottom
    (list (make-string w #\-)
          (current-prompt model)))
  
  (define log-space
    (- h (length bottom)))
  
  (append (visible-log log-space model) bottom))

(define (visible-log log-space model)
  (define log (model 'log))
  (define upscroll (model 'upscroll))
  
  (~> log
      (drop-right-upto upscroll)
      (take-right-upto log-space)
      (prefix-with-upto log-space "")))

(define (prefix-with-upto l n p)
  (~> (append (make-list n p) l)
      (take-right-upto n)))

(define (drop-right-upto l n)
  (~>> n
       (min (length l))
       (drop-right l)))

(define (take-right-upto l n)
  (~>> n
       (min (length l))
       (take-right l)))

(define (current-prompt model)
  (match (model 'options)

    [(list)
     "! No actions available."]

    [(focus _ (option title _ _ _) _)
     (string-append "> " title)]))

(define (exit? model)
  (model 'done))

(struct focus (above on below) #:transparent)

(define (focus-shift-up foc)
  (match foc
    [(focus (list) _ _) foc]
    [(focus (list h t ...) m b)
     (focus t h (cons m b))]))

(define (focus-shift-down foc)
  (match foc
    [(focus _ _ (list)) foc]
    [(focus a m (list h t ...))
     (focus (cons m a) h t)]))