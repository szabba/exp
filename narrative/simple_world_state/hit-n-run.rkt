#lang rackjure

(require "world-state.rkt")
(require "ui.rkt")

(provide run!)

;; TODO: some abstraction around "groups of functionality"?
;; basically I want modules of quality-choice packages
;; with some static checking ("this module never references an unnamed quality")
;; by way of limiting possible conditions and effects

;;; Entrypoint

(define (run!)
  (~>> (init) run-world!))

(define (init)
  (world-state init-state rules))

;;; Game definition variables

(define init-state {})
(define rules (set))

;;; Helpers

(define-syntax define-quality
  (syntax-rules ()
    ((_ name init)
     (begin
       (~>> init
            (init-state 'name)
            (set! init-state))
       (define name 'name)))))

(define-syntax define-qualities
  (syntax-rules ()
    ((_ (name init) ...)
     (begin
       (define-quality name init) ...))))

(define (quality? v)
  (dict-has-key? init-state v))

(define (symbol-prefix? v prefix)
  (~> v
      symbol->string
      (string-prefix? (~> prefix
                          symbol->string))))

(define (add-choice! opt)
  (~>> opt (set-add rules) (set! rules)))

(define (eff quality f)
  (λ args
    (let [(f-fed (apply f args))]
      (λ (qs)
        (~>> (qs quality #:else 0)
             f-fed
             (min 0)
             (qs quality))))))

;;; Health

(define-qualities
  (health/wounds 0)
  (health/eyes 2)
  (appearance/scars 0))

(define (eff/take-wounds n)
  (eff health/wounds #λ(- % n)))

(define eff/lose-eye
  (eff health/eyes #λ(- % 1)))

(define eff/scar
  (eff appearance/scars (λ (_) 1)))

;;; Travel

(define-qualities
  (can/travel 1))

(define (cond/can-travel-from? src)
  (λ (qs)
    (and (~> qs can/travel (= 1))
         (~> qs src (= 1)))))

(define (eff/travel src dst)
  (compose
   (eff src (λ (_) 0))
   (eff dst (λ (_) 1))))

(define (travel-route/barebone src dst)
  (travel-route
   (str "Travel from " src " to " dst)
   src
   dst
   (list (str "You arrive at " dst))))

(define (travel-route title src dst text
                      #:cond [cond/x? (λ (_) #t)]
                      #:eff [eff/x (list)])
  (choice
   title
   (λ (qs)
     (and ((cond/can-travel-from? src) qs)
          (cond/x? qs)))
   text
   (cons (eff/travel src dst) eff/x)))

;;; Regions and routes

(define-qualities
  (region/mezarka/clan-territory 1)
  (region/mezarka/city 0))

(for-each add-choice!
          (list (travel-route/barebone region/mezarka/clan-territory
                                       region/mezarka/city)
                (travel-route/barebone region/mezarka/city
                                       region/mezarka/clan-territory)))